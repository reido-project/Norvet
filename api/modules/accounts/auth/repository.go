package auth

import (
	"api/modules/accounts/_shared"

	"github.com/doug-martin/goqu/v9"
)

type Repository struct {
	db *goqu.Database
}

func NewRepository(db *goqu.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *_shared.User) error {
	result, err := r.db.
		Insert(_shared.TableName).
		Rows(goqu.Record{_shared.ColUsername: user.Username, _shared.ColEmail: user.Email, _shared.ColPassword: user.Password}).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = uint(id)

	return nil
}
