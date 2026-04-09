package model

type Login struct {
	username string `db:"username"`
	password string `db:"password"`
}
