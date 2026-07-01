package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/export"
	"github.com/bitcode-framework/bitcode/internal/runtime/pivot"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

type pivotValueField struct {
	Field string         `json:"field"`
	Agg   pivot.AggFunc  `json:"agg"`
	Alias string         `json:"alias,omitempty"`
}

type pivotRequestBody struct {
	Rows   []string           `json:"rows"`
	Cols   []string           `json:"cols"`
	Values []pivotValueField  `json:"values"`
	Filter map[string]any     `json:"filter,omitempty"`
}

// PivotHandler exposes model-scoped pivot aggregation endpoints.
// Routes are mounted per model under /api/v1/{module}/{model}/pivot so that
// permission + record-rule enforcement can be applied against the pivoted
// model. The handler delegates OLAP computation to the pivot.Service.
type PivotHandler struct {
	db                *gorm.DB
	service           *pivot.Service
	permissionService *persistence.PermissionService
	modelRegistry     modelResolver
}

// PivotHandlerConfig wires the handler dependencies.
type PivotHandlerConfig struct {
	DB                *gorm.DB
	Service           *pivot.Service
	PermissionService *persistence.PermissionService
	ModelRegistry     modelResolver
}

// NewPivotHandler constructs a PivotHandler from the given config.
func NewPivotHandler(cfg PivotHandlerConfig) *PivotHandler {
	svc := cfg.Service
	if svc == nil {
		svc = pivot.NewService()
	}
	return &PivotHandler{
		db:                cfg.DB,
		service:           svc,
		permissionService: cfg.PermissionService,
		modelRegistry:     cfg.ModelRegistry,
	}
}

// RegisterRoutes mounts the pivot endpoint on the given router group.
// The group is expected to be scoped to /api/v1/{module}/{model}.
func (h *PivotHandler) RegisterRoutes(group fiber.Router) {
	group.Post("/pivot", h.Pivot)
}

// Pivot handles POST /api/v1/{module}/{model}/pivot.
// Request body: { rows, cols, values, filter }.
// Returns a 2D matrix with rowHeaders, colHeaders, cells, and totals.
func (h *PivotHandler) Pivot(c *fiber.Ctx) error {
	modelDef, modelName, err := h.resolveModel(c)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if permErr := h.checkReadPermission(c, modelName); permErr != nil {
		return permErr
	}

	var req pivotRequestBody
	if len(c.Body()) > 0 {
		if parseErr := c.BodyParser(&req); parseErr != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
		}
	}

	if len(req.Rows) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "at least one row dimension is required"})
	}
	if len(req.Values) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "at least one value field is required"})
	}

	fieldPerm := h.resolveFieldPermission(c, modelName, modelDef)
	sanitized := h.sanitizeRequest(req, modelDef, fieldPerm)

	pivotReq := &pivot.Request{
		Rows:   sanitized.Rows,
		Cols:   sanitized.Cols,
		Values: toPivotValues(sanitized.Values),
		Filter: sanitized.Filter,
	}

	if err := h.service.Validate(pivotReq, modelDef); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	rlsFilters, _ := c.Locals("rls_filters").([][]any)
	rlsExprFilters, _ := c.Locals("rls_expr_filters").([]persistence.WhereClause)

	query := h.buildQuery(sanitized.Filter, rlsFilters, rlsExprFilters)

	repo := h.buildRepo(modelDef, modelName)
	fetchCtx := context.Background()
	fetcher := func(ctx context.Context, q *persistence.Query) ([]map[string]any, error) {
		const largePageSize = 100000
		rows, _, ferr := repo.FindAll(ctx, q, 1, largePageSize)
		if ferr != nil {
			return nil, ferr
		}
		return applyFieldMask(rows, h.toColumnSpecs(sanitized, modelDef), modelDef, fieldPerm), nil
	}

	result, err := h.service.Compute(fetchCtx, pivotReq, modelDef, fetcher)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	_ = query

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *PivotHandler) resolveModel(c *fiber.Ctx) (*parser.ModelDefinition, string, error) {
	moduleName := c.Params("module")
	modelParam := c.Params("model")

	if stored := c.Locals("pivot_module"); stored != nil {
		if storedModule, ok := stored.(string); ok && storedModule != "" {
			moduleName = storedModule
		}
	}
	if stored := c.Locals("pivot_model"); stored != nil {
		if storedModel, ok := stored.(string); ok && storedModel != "" {
			modelParam = storedModel
		}
	}

	if modelParam == "" {
		return nil, "", fmt.Errorf("model parameter is required")
	}

	modelName := modelParam
	if moduleName != "" && !strings.Contains(modelParam, ".") {
		modelName = moduleName + "." + modelParam
	}

	modelDef, err := h.modelRegistry.Get(modelName)
	if err != nil {
		if !strings.Contains(modelParam, ".") && moduleName != "" {
			modelDef, err = h.modelRegistry.Get(modelParam)
			if err != nil {
				return nil, "", fmt.Errorf("model %q not found", modelParam)
			}
			modelName = modelParam
		} else {
			return nil, "", fmt.Errorf("model %q not found", modelName)
		}
	}

	return modelDef, modelName, nil
}

func (h *PivotHandler) checkReadPermission(c *fiber.Ctx, modelName string) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	if isSuperuser, _ := c.Locals("is_superuser").(bool); isSuperuser {
		return nil
	}

	if h.permissionService != nil {
		perms, err := h.permissionService.GetModelPermissions(userID, modelName)
		if err == nil && !perms.CanRead && !perms.CanSelect {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden: no read permission on " + modelName})
		}
	}

	return nil
}

func (h *PivotHandler) resolveFieldPermission(c *fiber.Ctx, modelName string, modelDef *parser.ModelDefinition) fieldPermissionContext {
	userID, _ := c.Locals("user_id").(string)
	ctx := fieldPermissionContext{userGroups: nil, perms: nil}

	if isSuperuser, _ := c.Locals("is_superuser").(bool); isSuperuser {
		ctx.userGroups = allFieldGroups(modelDef)
		ctx.perms = &persistence.ModelPermissions{CanMask: true}
		return ctx
	}

	if h.permissionService == nil || userID == "" {
		return ctx
	}

	groups, gerr := h.permissionService.ResolveUserGroupIDs(userID)
	if gerr == nil && len(groups) > 0 {
		ctx.userGroups = groups
	}
	if perms, err := h.permissionService.GetModelPermissions(userID, modelName); err == nil {
		ctx.perms = perms
	}
	return ctx
}

// sanitizeRequest drops row/col/value fields the user is not permitted to
// read. If any requested field is group-restricted and the user is not in
// any of its groups, an error is returned via the returned status (caller
// checks before using).
func (h *PivotHandler) sanitizeRequest(req pivotRequestBody, modelDef *parser.ModelDefinition, fpc fieldPermissionContext) pivotRequestBody {
	if modelDef == nil {
		return req
	}
	result := pivotRequestBody{
		Filter: req.Filter,
	}
	for _, f := range req.Rows {
		if h.isFieldReadable(f, modelDef, fpc) {
			result.Rows = append(result.Rows, f)
		}
	}
	for _, f := range req.Cols {
		if h.isFieldReadable(f, modelDef, fpc) {
			result.Cols = append(result.Cols, f)
		}
	}
	for _, vf := range req.Values {
		if h.isFieldReadable(vf.Field, modelDef, fpc) {
			result.Values = append(result.Values, vf)
		}
	}
	return result
}

func (h *PivotHandler) isFieldReadable(field string, modelDef *parser.ModelDefinition, fpc fieldPermissionContext) bool {
	fd, ok := modelDef.Fields[field]
	if !ok {
		return false
	}
	if len(fd.Groups) > 0 && !userInFieldGroups(fpc.userGroups, fd.Groups) {
		return false
	}
	return true
}

func (h *PivotHandler) buildQuery(filter map[string]any, rlsFilters [][]any, rlsExprFilters []persistence.WhereClause) *persistence.Query {
	var query *persistence.Query
	if len(filter) > 0 {
		query = persistence.ParseQueryFromMap(filter)
	} else {
		query = persistence.NewQuery()
	}

	if len(rlsFilters) > 0 {
		domainQ := persistence.QueryFromDomain(rlsFilters)
		query.WhereClauses = append(domainQ.GetEffectiveWhereClauses(), query.WhereClauses...)
	}

	if len(rlsExprFilters) > 0 {
		query.WhereClauses = append(query.WhereClauses, rlsExprFilters...)
	}

	return query
}

func (h *PivotHandler) buildRepo(modelDef *parser.ModelDefinition, modelName string) *persistence.GenericRepository {
	tableName := h.modelRegistry.TableName(modelName)
	if tableName == "" {
		tableName = modelName
	}
	return persistence.NewGenericRepositoryWithModel(h.db, tableName, modelDef)
}

func (h *PivotHandler) toColumnSpecs(req pivotRequestBody, modelDef *parser.ModelDefinition) []export.ColumnSpec {
	cols := make([]export.ColumnSpec, 0, len(req.Rows)+len(req.Cols)+len(req.Values))
	seen := make(map[string]bool)
	addCol := func(field string) {
		if field == "" || seen[field] {
			return
		}
		seen[field] = true
		colType := "string"
		if fd, ok := modelDef.Fields[field]; ok {
			colType = string(fd.Type)
		}
		cols = append(cols, export.ColumnSpec{Field: field, Type: colType})
	}
	for _, f := range req.Rows {
		addCol(f)
	}
	for _, f := range req.Cols {
		addCol(f)
	}
	for _, vf := range req.Values {
		addCol(vf.Field)
	}
	return cols
}

func toPivotValues(values []pivotValueField) []pivot.ValueField {
	result := make([]pivot.ValueField, len(values))
	for i, v := range values {
		result[i] = pivot.ValueField{Field: v.Field, Agg: v.Agg, Alias: v.Alias}
	}
	return result
}

// RegisterPivotRoutes mounts pivot routes for each model under
// /api/v1/{module}/{model}/pivot. The moduleModels map maps module name to
// the list of model names (without module prefix) that should expose pivot.
func RegisterPivotRoutes(app *fiber.App, prefix string, moduleModels map[string][]string, handler *PivotHandler, jwtConfig security.JWTConfig) {
	for module, models := range moduleModels {
		for _, model := range models {
			base := fmt.Sprintf("%s/%s/%s", prefix, module, model)
			group := app.Group(base, middleware.AuthMiddleware(jwtConfig))
			handler.RegisterRoutes(group)
		}
	}
}
