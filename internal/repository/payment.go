package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
	*sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db}
}

func (r *PaymentRepository) UpdateBalance(ctx context.Context, accountID, amount int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	if _, err := r.ExecContext(ctx, query, amount, accountID); err != nil {
		return err
	}

	return nil
}
