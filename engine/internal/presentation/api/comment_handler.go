package api

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/comment"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

// CommentHandler exposes record-scoped comment endpoints. Comments are
// polymorphic (morph_to "commentable") so a single set of routes is mounted
// per module under /api/v1/{module}/comments. Permission + record-rule
// enforcement is performed against the commented model so that a user can
// only read/write comments on records they can already access.
type CommentHandler struct {
	service           *comment.Service
	permissionService *persistence.PermissionService
	recordRuleService *persistence.RecordRuleService
	modelRegistry     interface {
		Get(string) (*parser.ModelDefinition, error)
		MorphType(string) string
	}
}

// CommentHandlerConfig wires the handler dependencies.
type CommentHandlerConfig struct {
	Service           *comment.Service
	PermissionService *persistence.PermissionService
	RecordRuleService *persistence.RecordRuleService
	ModelRegistry     interface {
		Get(string) (*parser.ModelDefinition, error)
		MorphType(string) string
	}
}

// NewCommentHandler constructs a CommentHandler from the given config.
func NewCommentHandler(cfg CommentHandlerConfig) *CommentHandler {
	return &CommentHandler{
		service:           cfg.Service,
		permissionService: cfg.PermissionService,
		recordRuleService: cfg.RecordRuleService,
		modelRegistry:     cfg.ModelRegistry,
	}
}

// RegisterRoutes mounts the comment endpoints on the given router group.
func (h *CommentHandler) RegisterRoutes(group fiber.Router) {
	group.Get("/comments", h.List)
	group.Post("/comments", h.Create)
	group.Delete("/comments/:id", h.Delete)
}

// List returns the threaded comment list for a given record.
// GET /api/v1/{module}/comments?record_id=X&page=1&page_size=50
func (h *CommentHandler) List(c *fiber.Ctx) error {
	model := c.Params("module")
	recordID := c.Query("record_id")
	if recordID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "record_id is required"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size", "50"))
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	if err := h.checkPermission(c, model, recordID, "read"); err != nil {
		return err
	}

	morphType := h.resolveMorphType(model)

	records, total, err := h.service.ListByRecord(c.Context(), morphType, recordID, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	threaded := buildThread(records)
	return c.JSON(fiber.Map{
		"data":      threaded,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

type createCommentRequest struct {
	Body     string `json:"body"`
	RecordID string `json:"record_id"`
	ParentID string `json:"parent_id"`
}

// Create parses @mentions, creates the comment, fires notifications, and
// broadcasts a WebSocket event.
// POST /api/v1/{module}/comments  body: { record_id, body, parent_id? }
func (h *CommentHandler) Create(c *fiber.Ctx) error {
	model := c.Params("module")

	var req createCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if strings.TrimSpace(req.Body) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "body is required"})
	}
	if req.RecordID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "record_id is required"})
	}

	if err := h.checkPermission(c, model, req.RecordID, "write"); err != nil {
		return err
	}

	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	morphType := h.resolveMorphType(model)

	var parentID *string
	if req.ParentID != "" {
		parentID = &req.ParentID
	}

	created, err := h.service.Create(c.Context(), comment.CreateCommentInput{
		Body:            req.Body,
		CommentableType: morphType,
		CommentableID:   req.RecordID,
		ParentID:        parentID,
		AuthorID:        userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": toCommentJSON(*created)})
}

// Delete removes a comment. Only the author or a user with delete permission
// (or a superuser) may delete.
// DELETE /api/v1/{module}/comments/:id
func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "comment id is required"})
	}

	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	existing, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "comment not found"})
	}

	isAuthor := existing.AuthorID != nil && *existing.AuthorID == userID
	isSuperuser, _ := c.Locals("is_superuser").(bool)
	if !isAuthor && !isSuperuser {
		if h.permissionService != nil {
			perms, perr := h.permissionService.GetModelPermissions(userID, "comment")
			if perr != nil || !perms.CanDelete {
				return c.Status(403).JSON(fiber.Map{"error": "forbidden: only the author or a user with delete permission can delete a comment"})
			}
		} else {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden: only the author can delete a comment"})
		}
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "deleted"})
}

func (h *CommentHandler) resolveMorphType(modelName string) string {
	if h.modelRegistry != nil {
		return h.modelRegistry.MorphType(modelName)
	}
	return modelName
}

func (h *CommentHandler) checkPermission(c *fiber.Ctx, modelName, recordID, operation string) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return fiber.NewError(401, "authentication required")
	}

	if isSuperuser, _ := c.Locals("is_superuser").(bool); isSuperuser {
		return nil
	}

	if h.permissionService != nil {
		perms, err := h.permissionService.GetModelPermissions(userID, modelName)
		if err == nil {
			switch operation {
			case "read":
				if !perms.CanSelect {
					return fiber.NewError(403, "forbidden: no read permission on "+modelName)
				}
			case "write":
				if !perms.CanWrite && !perms.CanCreate {
					return fiber.NewError(403, "forbidden: no write permission on "+modelName)
				}
			}
		}
	}

	if h.recordRuleService != nil {
		allowed, err := h.recordRuleService.CheckRecordAccess(userID, modelName, recordID, operation)
		if err != nil {
			return fiber.NewError(403, "forbidden: record rule check failed")
		}
		if !allowed {
			return fiber.NewError(403, "forbidden: record rule denied access")
		}
	}

	return nil
}

// RegisterCommentRoutes mounts comment routes for each module under
// /api/v1/{module}/comments. Comments are polymorphic and available for any
// model; a single :module param group drives morph-type resolution. Only the
// modules listed in `modules` are accepted — other module values get 403.
func RegisterCommentRoutes(app *fiber.App, prefix string, modules []string, handler *CommentHandler, jwtConfig security.JWTConfig) {
	allowed := make(map[string]bool, len(modules))
	for _, m := range modules {
		allowed[m] = true
	}
	group := app.Group(prefix+"/:module", middleware.AuthMiddleware(jwtConfig))
	group.Use(func(c *fiber.Ctx) error {
		mod := c.Params("module")
		if !allowed[mod] {
			return c.Status(403).JSON(fiber.Map{"error": "comments not enabled for module: " + mod})
		}
		return c.Next()
	})
	handler.RegisterRoutes(group)
}

// toCommentJSON converts a CommentRecord to a JSON-friendly map.
func toCommentJSON(r comment.CommentRecord) fiber.Map {
	m := fiber.Map{
		"id":               r.ID,
		"body":             r.Body,
		"commentable_type": r.CommentableType,
		"commentable_id":   r.CommentableID,
		"created_at":       r.CreatedAt,
		"updated_at":       r.UpdatedAt,
	}
	if r.ParentID != nil {
		m["parent_id"] = *r.ParentID
	}
	if r.AuthorID != nil {
		m["author_id"] = *r.AuthorID
	}
	if len(r.Mentions) > 0 {
		m["mentions"] = r.Mentions
	}
	return m
}

// commentNode is a threaded comment with its replies.
type commentNode struct {
	comment.CommentRecord
	Replies []commentNode `json:"replies"`
}

// buildThread arranges a flat list of comments (ordered ASC by created_at)
// into a threaded structure using parent_id. Top-level comments (no parent_id)
// appear first; replies are nested under their parent.
func buildThread(records []comment.CommentRecord) []commentNode {
	byID := make(map[string]*commentNode, len(records))
	for _, r := range records {
		byID[r.ID] = &commentNode{CommentRecord: r}
	}

	var rootIDs []string
	childByParent := make(map[string][]string, len(records))
	for i := range records {
		r := &records[i]
		if r.ParentID == nil || *r.ParentID == "" {
			rootIDs = append(rootIDs, r.ID)
		} else if _, ok := byID[*r.ParentID]; ok {
			childByParent[*r.ParentID] = append(childByParent[*r.ParentID], r.ID)
		} else {
			rootIDs = append(rootIDs, r.ID)
		}
	}

	roots := make([]commentNode, 0, len(rootIDs))
	for _, id := range rootIDs {
		roots = append(roots, *buildNode(byID, childByParent, id))
	}
	return roots
}

func buildNode(byID map[string]*commentNode, childByParent map[string][]string, id string) *commentNode {
	node := byID[id]
	for _, childID := range childByParent[id] {
		node.Replies = append(node.Replies, *buildNode(byID, childByParent, childID))
	}
	return node
}
