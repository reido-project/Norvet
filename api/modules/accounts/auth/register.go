package auth

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v3"
)

func Register(route fiber.Router, db *goqu.Database) {
	// Feature
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Routes
	auth := route.Group("/auth")

	auth.Post("/register", handler.Register)
}
