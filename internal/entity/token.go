package entity

type Token struct {
	ID          int64  `db:"id"`
	TokenString string `db:"token_string"`
	IsValid     bool   `db:"is_valid"`
}
