package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/notification"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

// NotificationHandler exposes user-scoped notification endpoints behind
// authentication. All endpoints operate on the authenticated user's
// notifications — a user can never list, read, or mutate another user's
// notifications.
type NotificationHandler struct {
	service *notification.Service
}

// NotificationHandlerConfig wires the handler dependencies.
type NotificationHandlerConfig struct {
	Service *notification.Service
}

// NewNotificationHandler constructs a NotificationHandler from the given config.
func NewNotificationHandler(cfg NotificationHandlerConfig) *NotificationHandler {
	return &NotificationHandler{service: cfg.Service}
}

// RegisterRoutes mounts the notification endpoints on the given router group.
func (h *NotificationHandler) RegisterRoutes(group fiber.Router) {
	group.Get("/notifications", h.List)
	group.Post("/notifications/:id/read", h.MarkRead)
	group.Post("/notifications/read-all", h.MarkAllRead)
	group.Get("/notifications/unread-count", h.UnreadCount)
}

// List returns the authenticated user's notifications, unread-first, with
// pagination.
// GET /api/v1/notifications?page=1&page_size=20&unread_only=false
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	if pageSize < 1 || pageSize > 200 {
		pageSize = 20
	}
	unreadOnly := c.Query("unread_only") == "true" || c.Query("unread_only") == "1"

	result, err := h.service.List(c.Context(), notification.ListOptions{
		UserID:     userID,
		Page:       page,
		PageSize:   pageSize,
		UnreadOnly: unreadOnly,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// MarkRead marks a single notification as read.
// POST /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "notification id is required"})
	}

	if err := h.service.MarkRead(c.Context(), userID, id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "marked as read", "id": id})
}

// MarkAllRead marks all of the authenticated user's unread notifications as read.
// POST /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllRead(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	count, err := h.service.MarkAllRead(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "all marked as read", "count": count})
}

// UnreadCount returns the number of unread notifications for the authenticated user.
// GET /api/v1/notifications/unread-count
func (h *NotificationHandler) UnreadCount(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
	}

	count, err := h.service.UnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"count": count})
}

// RegisterNotificationRoutes mounts the notification endpoints at
// {prefix}/notifications behind AuthMiddleware.
func RegisterNotificationRoutes(app *fiber.App, prefix string, handler *NotificationHandler, jwtConfig security.JWTConfig) {
	group := app.Group(prefix, middleware.AuthMiddleware(jwtConfig))
	handler.RegisterRoutes(group)
}
