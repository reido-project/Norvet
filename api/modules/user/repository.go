package user

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type Repository struct {
	db *goqu.Database
}

func NewRepository(db *goqu.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) findOne(where goqu.Ex) (*User, error) {
	var user User

	found, err := r.db.
		From(tableName).
		Select(colID, colUsername, colEmail).
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

func (r *Repository) FindAll() ([]User, error) {
	var users []User

	err := r.db.
		From(tableName).
		Select(colID, colUsername, colEmail).
		ScanStructs(&users)

	return users, err
}

func (r *Repository) FindByID(id uint) (*User, error) {
	return r.findOne(goqu.Ex{colID: id})
}

func (r *Repository) FindByUsername(username string) (*User, error) {
	return r.findOne(goqu.Ex{colUsername: username})
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	return r.findOne(goqu.Ex{colEmail: email})
}

func (r *Repository) Create(user *User) error {
	result, err := r.db.
		Insert(tableName).
		Rows(goqu.Record{colUsername: user.Username, colEmail: user.Email, colPassword: user.Password}).
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

func (r *Repository) Update(user *User) error {
	result, err := r.db.
		Update(tableName).
		Set(goqu.Record{colUsername: user.Username, colEmail: user.Email}).
		Where(goqu.Ex{colID: user.ID}).
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
		Delete(tableName).
		Where(goqu.Ex{colID: id}).
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
