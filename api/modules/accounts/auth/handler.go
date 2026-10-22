package auth

import (
	"api/libs/validator"
	"api/modules/accounts/_shared"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(ctx *fiber.Ctx) error {
	var req _shared.CreateUserRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := validator.Validate(req); len(errs) > 0 {
		return ctx.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"errors":  errs,
		})
	}

	user, err := h.service.Register(&req)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to create accounts",
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}
