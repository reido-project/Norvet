package auth

import (
	"github.com/doug-martin/goqu/v9"
)

type repository struct {
	db *goqu.Database
}

func newRepository(db *goqu.Database) *repository {
	return &repository{db}
}

func (r *repository) create(user *userModel) error {
	result, err := r.db.
		Insert(tableName).
		Rows(goqu.Record{
			usernameCol:    user.Username,
			profileNameCol: user.ProfileName,
			emailCol:       user.Email,
			passwordCol:    user.Password,
			createdAtCol:   user.CreatedAt,
			updatedAtCol:   user.UpdatedAt,
		}).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint32(id)
	return nil
}
