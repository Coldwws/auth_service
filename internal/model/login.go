package model

type Login struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
