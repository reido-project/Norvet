package auth

import (
	"github.com/gofiber/fiber/v3"
	"gitlab.com/Norvet/nvterr"
	"gitlab.com/Norvet/nvtlog"
	"gitlab.com/Norvet/nvtval"
)

type handler struct {
	service *service
}

func newHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) register(ctx fiber.Ctx) error {
	var req registerRequest

	err := ctx.Bind().Body(&req)

	if err != nil {
		nerr := nvterr.ErrBadRequest
		return nvtlog.Response(ctx, nerr, fiber.Map{
			"success": false,
			"message": nerr.Message(),
		})
	}

	if errs := nvtval.SanitizeAndValidate(&req); len(errs) > 0 {
		return nvtlog.Response(ctx, nvterr.ErrValidationFailed, fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}

	response, err := h.service.register(&req)

	if err != nil {
		nerr := nvterr.From(err)
		return nvtlog.ErrorResponse(ctx, nerr, fiber.Map{
			"success": false,
			"message": nerr,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}
