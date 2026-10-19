package user

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db *goqu.Database) {
	// Feature
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Routes
	users := app.Group("/users")

	users.Get("/all", handler.GetAll)
	users.Get("/id", handler.GetByID)
	users.Get("/username", handler.GetByUsername)
	users.Get("/email", handler.GetByEmail)

	users.Post("/", handler.Create)
	users.Put("/:id", handler.Update)
	users.Delete("/:id", handler.Delete)
}
