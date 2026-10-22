package modules

import (
	"api/modules/accounts"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db *goqu.Database) {
	// Register your services here
	accounts.Register(app, db)
}
