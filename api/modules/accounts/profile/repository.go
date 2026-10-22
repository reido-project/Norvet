package profile

import (
	"api/modules/accounts/_shared"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type Repository struct {
	db *goqu.Database
}

func NewRepository(db *goqu.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) findOne(where goqu.Ex) (*_shared.User, error) {
	var user _shared.User

	found, err := r.db.
		From(_shared.TableName).
		Select(_shared.ColID, _shared.ColUsername, _shared.ColEmail).
		Where(where).
		ScanStruct(&user)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, sql.ErrNoRows
	}

	return &user, nil
}

func (r *Repository) FindAll() ([]_shared.User, error) {
	var users []_shared.User

	err := r.db.
		From(_shared.TableName).
		Select(_shared.ColID, _shared.ColUsername, _shared.ColEmail).
		ScanStructs(&users)

	return users, err
}

func (r *Repository) FindByID(id uint) (*_shared.User, error) {
	return r.findOne(goqu.Ex{_shared.ColID: id})
}

func (r *Repository) FindByUsername(username string) (*_shared.User, error) {
	return r.findOne(goqu.Ex{_shared.ColUsername: username})
}

func (r *Repository) FindByEmail(email string) (*_shared.User, error) {
	return r.findOne(goqu.Ex{_shared.ColEmail: email})
}

func (r *Repository) Update(user *_shared.User) error {
	result, err := r.db.
		Update(_shared.TableName).
		Set(goqu.Record{_shared.ColUsername: user.Username, _shared.ColEmail: user.Email}).
		Where(goqu.Ex{_shared.ColID: user.ID}).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(id uint) error {
	result, err := r.db.
		Delete(_shared.TableName).
		Where(goqu.Ex{_shared.ColID: id}).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
