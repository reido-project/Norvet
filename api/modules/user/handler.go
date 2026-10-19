package user

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.service.GetAll()

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to get users",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

func (h *Handler) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Query("id"), 10, 64)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	user, err := h.service.GetByID(uint(id))

	err = errorCheck(err, ctx)

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) GetByUsername(ctx *fiber.Ctx) error {
	username := strings.TrimSpace(ctx.Query("username"))

	if username == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "username cannot be empty",
		})
	}

	user, err := h.service.GetByUsername(username)

	err = errorCheck(err, ctx)

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) GetByEmail(ctx *fiber.Ctx) error {
	email := strings.TrimSpace(ctx.Query("email"))

	if email == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "email cannot be empty",
		})
	}

	user, err := h.service.GetByEmail(email)

	err = errorCheck(err, ctx)

	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var req CreateUserRequest
	err := ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return ctx.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "username, email, and password are required",
		})
	}

	user, err := h.service.Create(&req)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to create user",
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	var req UpdateUserRequest
	err = ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" {
		return ctx.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "username and email cannot be empty",
		})
	}

	user, err := h.service.Update(uint(id), &req)

	if errors.Is(err, sql.ErrNoRows) {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to update user",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	err = h.service.Delete(uint(id))

	if errors.Is(err, sql.ErrNoRows) {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to delete user",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "user has been deleted",
	})
}

// ============ Helper ============ //
func errorCheck(err error, ctx *fiber.Ctx) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "failed to get user",
		})
	}

	return nil
}
