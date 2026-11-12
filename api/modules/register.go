package modules

import (
	"api/modules/auth"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App, db *goqu.Database) {
	auth.Register(app, db)
}
