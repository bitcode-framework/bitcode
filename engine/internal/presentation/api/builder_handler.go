package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/builder"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

type BuilderHandler struct {
	viewValidator *builder.ViewValidator
}

type BuilderHandlerConfig struct {
	ViewValidator *builder.ViewValidator
}

func NewBuilderHandler(cfg BuilderHandlerConfig) *BuilderHandler {
	return &BuilderHandler{
		viewValidator: cfg.ViewValidator,
	}
}

func (h *BuilderHandler) ValidateView(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"valid":  false,
			"errors": []builder.ValidationIssue{{Field: "root", Message: "request body is required"}},
		})
	}

	result := h.viewValidator.Validate(body)
	return c.JSON(result)
}

func RegisterBuilderRoutes(app *fiber.App, prefix string, handler *BuilderHandler, jwtConfig security.JWTConfig) {
	group := app.Group(prefix, middleware.AuthMiddleware(jwtConfig))
	group.Post("/views/validate", handler.ValidateView)
}


