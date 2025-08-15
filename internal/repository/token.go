package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

type TokenRepository struct {
	*sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db}
}

func (r *TokenRepository) CreateToken(ctx context.Context, token *entity.Token) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// question: pass values from service layer
	// or delegate it to database default values

	query := `INSERT INTO tokens (id, token_string) VALUES ($1, $2)`
	if _, err := r.ExecContext(ctx, query, token.UserID, token.TokenString); err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) GetToken(ctx context.Context, tokenString string) (*entity.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var token entity.Token
	query := `SELECT * FROM tokens WHERE token_string = $1`
	if err := r.QueryRowxContext(ctx, query, tokenString).StructScan(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) InvalidateUserTokens(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE tokens SET is_valid = FALSE WHERE id = $1`
	if _, err := r.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
