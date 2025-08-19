package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

type RentRepository struct {
	*sqlx.DB
}

func NewRentRepository(db *sqlx.DB) *RentRepository {
	return &RentRepository{db}
}

func (r *RentRepository) StartRent(ctx context.Context, rent *entity.Rent) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var id int64
	query := `INSERT INTO rents (transport_id, user_id, time_start, time_end, price_of_unit, price_type, final_price) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	if err := r.QueryRowxContext(ctx, query, rent.TransportID, rent.UserID, rent.TimeStart, rent.TimeEnd, rent.PriceOfUnit, rent.PriceType, rent.FinalPrice).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RentRepository) EndRent(ctx context.Context, id int64, lat, long float64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	tx, err := r.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var rentID int64
	query := `
		UPDATE rents
		SET time_end = NOW(),
		    final_price = CASE price_type
		        WHEN 'Minutes' THEN CEIL(EXTRACT(EPOCH FROM (NOW() - time_start)) / 60) * price_of_unit
		        ELSE CEIL(EXTRACT(EPOCH FROM (NOW() - time_start)) / 86400) * price_of_unit
		    END
		WHERE id = $1
		RETURNING id
	`
	if err := tx.QueryRowxContext(ctx, query, id).Scan(&rentID); err != nil {
		return fmt.Errorf("update rent: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit rent: %w", err)
	}

	return nil
}

func (r *RentRepository) GetByID(ctx context.Context, id int64) (*entity.Rent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var rent entity.Rent
	query := `SELECT * FROM rents WHERE id = $1`
	if err := r.QueryRowxContext(ctx, query, id).StructScan(&rent); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &rent, nil
}

func (r *RentRepository) GetHistoryByUser(ctx context.Context, userID int64) ([]entity.Rent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	rents := make([]entity.Rent, 0, 100)
	query := `SELECT * FROM rents WHERE user_id = $1`
	rows, err := r.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rent entity.Rent
		if err := rows.StructScan(&rent); err != nil {
			return nil, err
		}

		rents = append(rents, rent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rents, nil
}

func (r *RentRepository) GetHistoryByTransport(ctx context.Context, transportID int64) ([]entity.Rent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	rents := make([]entity.Rent, 0, 100)
	query := `SELECT * FROM rents WHERE transport_id = $1`
	rows, err := r.QueryxContext(ctx, query, transportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rent entity.Rent
		if err := rows.StructScan(&rent); err != nil {
			return nil, err
		}

		rents = append(rents, rent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rents, nil
}

func (r *RentRepository) ListActive(ctx context.Context) ([]entity.Rent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	rents := make([]entity.Rent, 0, 100)
	query := `SELECT * FROM rents WHERE time_end IS NULL`
	rows, err := r.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rent entity.Rent
		if err := rows.StructScan(&rent); err != nil {
			return nil, err
		}

		rents = append(rents, rent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rents, nil
}

func (r *RentRepository) Update(ctx context.Context, rent *entity.Rent) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE rents SET transport_id = $1, user_id = $2, time_start = $3, time_end = $4, price_of_unit = $5, price_type = $6, final_price = $7 WHERE id = $8`
	if _, err := r.ExecContext(ctx, query, rent.TransportID, rent.UserID, rent.TimeStart, rent.TimeEnd, rent.PriceOfUnit, rent.PriceType, rent.FinalPrice, rent.ID); err != nil {
		return err
	}

	return nil
}

func (r *RentRepository) UpdateLastBilledTime(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE rents SET last_billed_at = NOW() WHERE id = $1`
	if _, err := r.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (r *RentRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `DELETE FROM rents WHERE id = $1`
	if _, err := r.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
