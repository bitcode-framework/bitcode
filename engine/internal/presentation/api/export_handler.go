package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/export"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

const exportBatchSize = 1000

type columnRequest struct {
	Field string `json:"field"`
	Label string `json:"label"`
}

type sortRequest struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type exportRequestBody struct {
	Format  string          `json:"format"`
	Columns []columnRequest `json:"columns"`
	Filter  map[string]any  `json:"filter"`
	Sort    *sortRequest    `json:"sort"`
}

type modelResolver interface {
	Get(string) (*parser.ModelDefinition, error)
	TableName(string) string
}

// ExportHandler exposes model-scoped data export endpoints (CSV/XLSX/JSON).
// Routes are mounted per model under /api/v1/{module}/{model}/export so that
// permission + record-rule enforcement can be applied against the exported
// model. The handler streams rows in batches to avoid loading the entire
// dataset into memory.
type ExportHandler struct {
	db                *gorm.DB
	permissionService *persistence.PermissionService
	modelRegistry     modelResolver
}

// ExportHandlerConfig wires the handler dependencies.
type ExportHandlerConfig struct {
	DB                *gorm.DB
	PermissionService *persistence.PermissionService
	ModelRegistry     modelResolver
}

// NewExportHandler constructs an ExportHandler from the given config.
func NewExportHandler(cfg ExportHandlerConfig) *ExportHandler {
	return &ExportHandler{
		db:                cfg.DB,
		permissionService: cfg.PermissionService,
		modelRegistry:     cfg.ModelRegistry,
	}
}

// RegisterRoutes mounts POST + GET export endpoints on the given router group.
// The group is expected to be scoped to /api/v1/{module}/{model}.
func (h *ExportHandler) RegisterRoutes(group fiber.Router) {
	group.Post("/export", h.Export)
	group.Get("/export", h.ExportGet)
}

// Export handles POST /api/v1/{module}/{model}/export.
// Request body: { format, columns, filter, sort }.
func (h *ExportHandler) Export(c *fiber.Ctx) error {
	modelDef, modelName, err := h.resolveModel(c)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if permErr := h.checkExportPermission(c, modelName); permErr != nil {
		return permErr
	}

	var req exportRequestBody
	if len(c.Body()) > 0 {
		if parseErr := c.BodyParser(&req); parseErr != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
		}
	}

	format := strings.ToLower(req.Format)
	if format == "" {
		format = "csv"
	}
	if format != "csv" && format != "xlsx" && format != "json" {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("unsupported format %q (supported: csv, xlsx, json)", format)})
	}

	rlsFilters, _ := c.Locals("rls_filters").([][]any)
	rlsExprFilters, _ := c.Locals("rls_expr_filters").([]persistence.WhereClause)

	query := h.buildQuery(req.Filter, req.Sort, rlsFilters, rlsExprFilters)

	fieldPerm := h.resolveFieldPermission(c, modelName, modelDef)
	columns := h.sanitizeColumns(h.resolveColumns(modelDef, req.Columns), modelDef, fieldPerm)
	if len(columns) == 0 {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: no exportable columns permitted for this user"})
	}

	repo := h.buildRepo(modelDef, modelName)
	fetchCtx := context.Background()
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		rows, total, ferr := repo.FindAll(fetchCtx, query, page, pageSize)
		if ferr != nil {
			return nil, total, ferr
		}
		return applyFieldMask(rows, columns, modelDef, fieldPerm), total, nil
	}

	contentType, filename := exportMeta(modelName, format)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return h.streamExport(c, format, modelName, columns, fetcher)
}

// ExportGet handles GET /api/v1/{module}/{model}/export for browser-initiated
// downloads. Query params: format, columns (comma-separated field list),
// oql (filter), sort_field, sort_order.
func (h *ExportHandler) ExportGet(c *fiber.Ctx) error {
	modelDef, modelName, err := h.resolveModel(c)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if permErr := h.checkExportPermission(c, modelName); permErr != nil {
		return permErr
	}

	format := strings.ToLower(c.Query("format", "csv"))
	if format != "csv" && format != "xlsx" && format != "json" {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("unsupported format %q (supported: csv, xlsx, json)", format)})
	}

	var requested []columnRequest
	if colsParam := c.Query("columns"); colsParam != "" {
		for _, field := range strings.Split(colsParam, ",") {
			field = strings.TrimSpace(field)
			if field != "" {
				requested = append(requested, columnRequest{Field: field})
			}
		}
	}

	var sort *sortRequest
	if sortField := c.Query("sort_field"); sortField != "" {
		sort = &sortRequest{Field: sortField, Order: c.Query("sort_order", "asc")}
	}

	var filter map[string]any
	if oql := c.Query("oql"); oql != "" {
		pq, _, perr := persistence.ParseOQL(oql)
		if perr != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("invalid OQL: %s", perr.Error())})
		}
		filter = queryToMap(pq)
	}

	rlsFilters, _ := c.Locals("rls_filters").([][]any)
	rlsExprFilters, _ := c.Locals("rls_expr_filters").([]persistence.WhereClause)

	query := h.buildQuery(filter, sort, rlsFilters, rlsExprFilters)

	fieldPerm := h.resolveFieldPermission(c, modelName, modelDef)
	columns := h.sanitizeColumns(h.resolveColumns(modelDef, requested), modelDef, fieldPerm)
	if len(columns) == 0 {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: no exportable columns permitted for this user"})
	}

	repo := h.buildRepo(modelDef, modelName)
	fetchCtx := context.Background()
	fetcher := func(page, pageSize int) ([]map[string]any, int64, error) {
		rows, total, ferr := repo.FindAll(fetchCtx, query, page, pageSize)
		if ferr != nil {
			return nil, total, ferr
		}
		return applyFieldMask(rows, columns, modelDef, fieldPerm), total, nil
	}

	contentType, filename := exportMeta(modelName, format)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return h.streamExport(c, format, modelName, columns, fetcher)
}

func (h *ExportHandler) resolveModel(c *fiber.Ctx) (*parser.ModelDefinition, string, error) {
	moduleName := c.Params("module")
	modelParam := c.Params("model")

	if stored := c.Locals("export_module"); stored != nil {
		if storedModule, ok := stored.(string); ok && storedModule != "" {
			moduleName = storedModule
		}
	}
	if stored := c.Locals("export_model"); stored != nil {
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

func (h *ExportHandler) checkExportPermission(c *fiber.Ctx, modelName string) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	if isSuperuser, _ := c.Locals("is_superuser").(bool); isSuperuser {
		return nil
	}

	if h.permissionService != nil {
		perms, err := h.permissionService.GetModelPermissions(userID, modelName)
		if err == nil && !perms.CanExport {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden: no export permission on " + modelName})
		}
	}

	return nil
}

func (h *ExportHandler) buildQuery(filter map[string]any, sort *sortRequest, rlsFilters [][]any, rlsExprFilters []persistence.WhereClause) *persistence.Query {
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

	if sort != nil && sort.Field != "" {
		order := strings.ToLower(sort.Order)
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		query.Order(sort.Field, order)
	}

	return query
}

func (h *ExportHandler) resolveColumns(modelDef *parser.ModelDefinition, requested []columnRequest) []export.ColumnSpec {
	if len(requested) > 0 {
		cols := make([]export.ColumnSpec, 0, len(requested))
		for _, rc := range requested {
			label := rc.Label
			if label == "" {
				if fd, ok := modelDef.Fields[rc.Field]; ok && fd.Label != "" {
					label = fd.Label
				} else {
					label = rc.Field
				}
			}
			colType := "string"
			if fd, ok := modelDef.Fields[rc.Field]; ok {
				colType = string(fd.Type)
			}
			cols = append(cols, export.ColumnSpec{Field: rc.Field, Label: label, Type: colType})
		}
		return cols
	}

	return export.DefaultColumnsFromModel(modelDef)
}

func (h *ExportHandler) buildRepo(modelDef *parser.ModelDefinition, modelName string) *persistence.GenericRepository {
	tableName := h.modelRegistry.TableName(modelName)
	if tableName == "" {
		tableName = modelName
	}
	return persistence.NewGenericRepositoryWithModel(h.db, tableName, modelDef)
}

func queryToMap(q *persistence.Query) map[string]any {
	if q == nil {
		return nil
	}
	data, err := json.Marshal(q)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	return m
}

func exportMeta(modelName, format string) (string, string) {
	contentType := "text/csv"
	switch format {
	case "xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "json":
		contentType = "application/json"
	}
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s", strings.ReplaceAll(modelName, ".", "_"), timestamp, format)
	return contentType, filename
}

func (h *ExportHandler) streamExport(c *fiber.Ctx, format, modelName string, columns []export.ColumnSpec, fetcher export.FetchFunc) error {
	if format == "json" {
		return h.streamJSON(c, columns, fetcher)
	}
	if format == "csv" {
		return h.streamCSV(c, columns, fetcher)
	}
	return h.streamXLSX(c, modelName, columns, fetcher)
}

func (h *ExportHandler) streamCSV(c *fiber.Ctx, columns []export.ColumnSpec, fetcher export.FetchFunc) error {
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer w.Flush()
		w.Write([]byte{0xEF, 0xBB, 0xBF})

		labels := make([]string, len(columns))
		for i, col := range columns {
			labels[i] = col.Label
		}
		writeCSVRow(w, labels)

		page := 1
		for {
			rows, _, err := fetcher(page, exportBatchSize)
			if err != nil {
				return
			}
			if len(rows) == 0 {
				break
			}
			for _, row := range rows {
				values := make([]string, len(columns))
				for i, col := range columns {
					values[i] = formatCellValue(row[col.Field])
				}
				writeCSVRow(w, values)
			}
			if len(rows) < exportBatchSize {
				break
			}
			page++
		}
	})

	return nil
}

func (h *ExportHandler) streamJSON(c *fiber.Ctx, columns []export.ColumnSpec, fetcher export.FetchFunc) error {
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer w.Flush()
		w.WriteString("[")

		first := true
		page := 1
		for {
			rows, _, err := fetcher(page, exportBatchSize)
			if err != nil {
				return
			}
			if len(rows) == 0 {
				break
			}
			for _, row := range rows {
				obj := make(map[string]any, len(columns))
				for _, col := range columns {
					obj[col.Label] = row[col.Field]
				}
				if !first {
					w.WriteString(",")
				}
				first = false
				enc := json.NewEncoder(w)
				enc.SetEscapeHTML(false)
				enc.Encode(obj)
			}
			if len(rows) < exportBatchSize {
				break
			}
			page++
		}
		w.WriteString("]")
	})

	return nil
}

func (h *ExportHandler) streamXLSX(c *fiber.Ctx, modelName string, columns []export.ColumnSpec, fetcher export.FetchFunc) error {
	buf, err := export.GenerateXLSX(modelName, columns, fetcher, exportBatchSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Send(buf)
}

func writeCSVRow(w io.Writer, values []string) {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = strconv.Quote(v)
	}
	fmt.Fprintln(w, strings.Join(quoted, ","))
}

func formatCellValue(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// fieldPermissionContext bundles the per-user field-level permission data
// needed to sanitize export columns and mask sensitive values. It mirrors the
// field_filter.go helpers used by the CRUD handler so export respects the same
// field-level groups + mask rules as the regular API responses.
type fieldPermissionContext struct {
	userGroups []string
	perms      *persistence.ModelPermissions
}

// resolveFieldPermission builds the field-permission context for the
// authenticated user against the given model. Superusers get an unrestricted
// context (all groups, can_mask=true). When the permission service is nil or
// the user lookup fails, an empty (deny-sensitive) context is returned so the
// handler fails closed.
func (h *ExportHandler) resolveFieldPermission(c *fiber.Ctx, modelName string, modelDef *parser.ModelDefinition) fieldPermissionContext {
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

// allFieldGroups returns every group referenced by any field on the model.
// Used to grant superusers visibility on all group-restricted fields.
func allFieldGroups(modelDef *parser.ModelDefinition) []string {
	if modelDef == nil {
		return nil
	}
	seen := make(map[string]bool)
	var groups []string
	for _, fd := range modelDef.Fields {
		for _, g := range fd.Groups {
			if !seen[g] {
				seen[g] = true
				groups = append(groups, g)
			}
		}
	}
	return groups
}

// sanitizeColumns drops columns the user is not permitted to see. A field is
// dropped when it declares field-level `groups` and the user belongs to none
// of those groups. Masked fields are kept (the value is masked later in
// applyFieldMask). Unknown fields (not in the model definition) are kept for
// backward compatibility with custom columns.
func (h *ExportHandler) sanitizeColumns(columns []export.ColumnSpec, modelDef *parser.ModelDefinition, fpc fieldPermissionContext) []export.ColumnSpec {
	if modelDef == nil {
		return columns
	}
	result := make([]export.ColumnSpec, 0, len(columns))
	for _, col := range columns {
		fd, ok := modelDef.Fields[col.Field]
		if !ok {
			result = append(result, col)
			continue
		}
		if len(fd.Groups) > 0 && !userInFieldGroups(fpc.userGroups, fd.Groups) {
			continue
		}
		result = append(result, col)
	}
	return result
}

// applyFieldMask masks values of masked fields when the user lacks the
// `can_mask` permission. Non-masked fields pass through unchanged. This runs
// per batch inside the fetcher so masking is applied before any value reaches
// the format generator.
func applyFieldMask(rows []map[string]any, columns []export.ColumnSpec, modelDef *parser.ModelDefinition, fpc fieldPermissionContext) []map[string]any {
	if modelDef == nil || len(rows) == 0 {
		return rows
	}
	canMask := fpc.perms != nil && fpc.perms.CanMask
	if canMask {
		return rows
	}
	maskFields := make(map[string]parser.FieldDefinition)
	for _, col := range columns {
		if fd, ok := modelDef.Fields[col.Field]; ok && fd.Mask {
			maskFields[col.Field] = fd
		}
	}
	if len(maskFields) == 0 {
		return rows
	}
	for i := range rows {
		for field, fd := range maskFields {
			if v, ok := rows[i][field]; ok {
				rows[i][field] = maskValue(v, fd.MaskLength)
			}
		}
	}
	return rows
}

// RegisterExportRoutes mounts export routes for each model under
// /api/v1/{module}/{model}/export. The moduleModels map maps module name to
// the list of model names (without module prefix) that should expose export.
func RegisterExportRoutes(app *fiber.App, prefix string, moduleModels map[string][]string, handler *ExportHandler, jwtConfig security.JWTConfig) {
	for module, models := range moduleModels {
		for _, model := range models {
			base := fmt.Sprintf("%s/%s/%s", prefix, module, model)
			group := app.Group(base, middleware.AuthMiddleware(jwtConfig))
			handler.RegisterRoutes(group)
		}
	}
}
