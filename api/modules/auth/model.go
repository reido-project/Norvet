package auth

import "time"

const (
	tableName      = "users"
	usernameCol    = "username"
	profileNameCol = "profile_name"
	passwordCol    = "password"
	emailCol       = "email"
	createdAtCol   = "created_at"
	updatedAtCol   = "updated_at"
)

type userModel struct {
	ID          uint32    `db:"id"`
	Username    string    `db:"username"`
	ProfileName string    `db:"profile_name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type registerRequest struct {
	Username    string `json:"username" validate:"required,alphanum,min=3,max=20" san:"trim,lower"`
	ProfileName string `json:"profile_name" validate:"required,min=3,max=50" san:"trim"`
	Email       string `json:"email" validate:"required,email,max=255" san:"trim,lower"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

type registerResponse struct {
	ID          uint32    `json:"id"`
	Username    string    `json:"username"`
	ProfileName string    `json:"profile_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

type loginRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3,max=255" san:"trim,lower"`
	Password   string `json:"password" validate:"required"`
}
