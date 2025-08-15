package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

type AccountRepository struct {
	*sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// question: pass time values from service layer
	// or delegate it to database default values

	var id int64
	query := `INSERT INTO accounts (username, password_hash, is_admin, balance) VALUES ($1, $2, $3, $4) RETURNING id`
	if err := r.QueryRowContext(ctx, query, account.Username, account.PasswordHash, account.IsAdmin, account.Balance).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var account entity.Account
	query := `SELECT * FROM accounts WHERE id = $1`
	if err := r.QueryRowxContext(ctx, query, id).StructScan(account); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var account entity.Account
	query := `SELECT * FROM accounts WHERE username = $1`
	if err := r.QueryRowxContext(ctx, query, username).StructScan(account); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetByUsernameAndPassword(ctx context.Context, username, passwordHash string) (*entity.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var account entity.Account
	query := `SELECT * FROM accounts WHERE username = $1 AND password_hash = $2`
	if err := r.QueryRowxContext(ctx, query, username, passwordHash).StructScan(account); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) List(ctx context.Context, count, start int) ([]entity.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	accounts := make([]entity.Account, 0, 100)
	query := `SELECT * FROM accounts LIMIT $1 OFFSET $2`
	rows, err := r.QueryxContext(ctx, query, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account entity.Account
		if err := rows.StructScan(&account); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE accounts SET username = $1, password_hash = $2, is_admin = $3, balance = $4, updated_at = $5 WHERE id = $6`
	if _, err := r.ExecContext(ctx, query, account.Username, account.PasswordHash, account.IsAdmin, account.Balance, account.UpdatedAt, account.ID); err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `DELETE FROM accounts WHERE id = $1`
	if _, err := r.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
