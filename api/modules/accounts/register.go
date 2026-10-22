package accounts

import (
	"api/modules/accounts/auth"
	"api/modules/accounts/profile"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db *goqu.Database) {
	accounts := app.Group("/accounts")

	auth.Register(accounts, db)
	profile.Register(accounts, db)
}
