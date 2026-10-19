package user

const (
	tableName   = "users"
	colID       = "id"
	colUsername = "username"
	colEmail    = "email"
	colPassword = "password"
)

type User struct {
	ID       uint   `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
