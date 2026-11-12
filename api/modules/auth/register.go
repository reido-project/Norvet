package auth

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, db *goqu.Database) {
	// Dependency Injection
	repository := newRepository(db)
	service := newService(repository)
	handler := newHandler(service)

	// Routes
	auth := app.Group("/auth")

	auth.Post("/register", handler.register)
}
