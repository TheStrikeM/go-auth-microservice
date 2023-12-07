package models

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	PassHash string `db:"password"`
}
