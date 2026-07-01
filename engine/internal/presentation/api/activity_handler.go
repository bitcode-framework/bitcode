package api

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/bitcode-framework/bitcode/internal/compiler/parser"
	"github.com/bitcode-framework/bitcode/internal/infrastructure/persistence"
	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/activity"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

// ActivityHandler exposes the aggregated activity feed endpoint. When a
// model + record_id pair is supplied the handler enforces read permission +
// record-rule access on the target model (same pattern as CommentHandler),
// otherwise it returns the global timeline for the authenticated user.
type ActivityHandler struct {
	service           *activity.Service
	permissionService *persistence.PermissionService
	recordRuleService *persistence.RecordRuleService
	modelRegistry     interface {
		Get(string) (*parser.ModelDefinition, error)
		MorphType(string) string
	}
}

// ActivityHandlerConfig wires the handler dependencies.
type ActivityHandlerConfig struct {
	Service           *activity.Service
	PermissionService *persistence.PermissionService
	RecordRuleService *persistence.RecordRuleService
	ModelRegistry     interface {
		Get(string) (*parser.ModelDefinition, error)
		MorphType(string) string
	}
}

// NewActivityHandler constructs an ActivityHandler from the given config.
func NewActivityHandler(cfg ActivityHandlerConfig) *ActivityHandler {
	return &ActivityHandler{
		service:           cfg.Service,
		permissionService: cfg.PermissionService,
		recordRuleService: cfg.RecordRuleService,
		modelRegistry:     cfg.ModelRegistry,
	}
}

// RegisterRoutes mounts the activity endpoints on the given router group.
func (h *ActivityHandler) RegisterRoutes(group fiber.Router) {
	group.Get("/activity", h.List)
}

// List returns aggregated activity entries.
// GET /api/v1/activity?model=X&record_id=Y&user=Z&since=ISO8601&limit=N&page=P
func (h *ActivityHandler) List(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	model := c.Query("model")
	recordID := c.Query("record_id")
	userFilter := c.Query("user")

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var since time.Time
	if raw := c.Query("since"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err == nil {
			since = parsed
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "invalid 'since' timestamp (expected ISO8601)"})
		}
	}

	if model != "" && recordID != "" {
		if err := h.checkPermission(c, model, recordID, "read"); err != nil {
			return err
		}
	}

	entries, total, err := h.service.ListActivity(c.Context(), activity.ActivityQueryOptions{
		Model:    model,
		RecordID: recordID,
		UserID:   userFilter,
		Since:    since,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":      entries,
		"total":     total,
		"page":      page,
		"page_size": limit,
	})
}

func (h *ActivityHandler) checkPermission(c *fiber.Ctx, modelName, recordID, operation string) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
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
					return c.Status(403).JSON(fiber.Map{"error": "forbidden: no read permission on " + modelName})
				}
			case "write":
				if !perms.CanWrite && !perms.CanCreate {
					return c.Status(403).JSON(fiber.Map{"error": "forbidden: no write permission on " + modelName})
				}
			}
		}
	}

	if h.recordRuleService != nil {
		allowed, err := h.recordRuleService.CheckRecordAccess(userID, modelName, recordID, operation)
		if err == nil && !allowed {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden: record rule denied access"})
		}
	}

	return nil
}

// RegisterActivityRoutes mounts the activity feed endpoint at {prefix}/activity
// behind AuthMiddleware. A single global route serves both record-scoped and
// global-timeline queries (the scope is driven by query params).
func RegisterActivityRoutes(app *fiber.App, prefix string, handler *ActivityHandler, jwtConfig security.JWTConfig) {
	group := app.Group(prefix, middleware.AuthMiddleware(jwtConfig))
	handler.RegisterRoutes(group)
}
