package dto

type UpdateUserDto struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
