package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bitcode-framework/bitcode/internal/presentation/middleware"
	"github.com/bitcode-framework/bitcode/internal/runtime/builder"
	"github.com/bitcode-framework/bitcode/pkg/security"
)

type BuilderHandler struct {
	viewValidator    *builder.ViewValidator
	processValidator *builder.ProcessValidator
	modelValidator   *builder.ModelValidator
}

type BuilderHandlerConfig struct {
	ViewValidator    *builder.ViewValidator
	ProcessValidator *builder.ProcessValidator
	ModelValidator   *builder.ModelValidator
}

func NewBuilderHandler(cfg BuilderHandlerConfig) *BuilderHandler {
	return &BuilderHandler{
		viewValidator:    cfg.ViewValidator,
		processValidator: cfg.ProcessValidator,
		modelValidator:   cfg.ModelValidator,
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

func (h *BuilderHandler) ValidateProcess(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"valid":  false,
			"errors": []builder.ValidationIssue{{Field: "root", Message: "request body is required"}},
		})
	}

	result := h.processValidator.Validate(body)
	return c.JSON(result)
}

func (h *BuilderHandler) ValidateModel(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"valid":  false,
			"errors": []builder.ValidationIssue{{Field: "root", Message: "request body is required"}},
		})
	}

	result := h.modelValidator.Validate(body)
	return c.JSON(result)
}

func RegisterBuilderRoutes(app *fiber.App, prefix string, handler *BuilderHandler, jwtConfig security.JWTConfig) {
	group := app.Group(prefix, middleware.AuthMiddleware(jwtConfig))
	group.Post("/views/validate", handler.ValidateView)
	if handler.processValidator != nil {
		group.Post("/processes/validate", handler.ValidateProcess)
	}
	if handler.modelValidator != nil {
		group.Post("/models/validate", handler.ValidateModel)
	}
}


