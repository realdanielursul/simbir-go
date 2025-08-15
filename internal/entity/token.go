package entity

type Token struct {
	UserID      int64  `db:"user_id"`
	TokenString string `db:"token_string"`
	IsValid     bool   `db:"is_valid"`
}
