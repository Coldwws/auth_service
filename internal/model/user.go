package model

type UserInfo struct {
	Email string
	Role  string
}

type UserDB struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	Role         string `db:"role"`
	PasswordHash string `db:"password_hash"`
}

type Login struct {
	Email    string
	Password string
}
