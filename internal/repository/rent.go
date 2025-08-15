package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	tx, err := r.BeginTxx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var id int64
	query := `INSERT INTO rents (transport_id, user_id, time_start, time_end, price_of_unit, price_type, final_price) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	if err := tx.QueryRowxContext(ctx, query, rent.TransportID, rent.UserID, rent.TimeStart, rent.TimeEnd, rent.PriceOfUnit, rent.PriceType, rent.FinalPrice).Scan(&id); err != nil {
		return -1, err
	}

	query = `UPDATE transports SET can_be_rented = FALSE WHERE id = $1`
	if _, err := tx.ExecContext(ctx, query, rent.TransportID); err != nil {
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *RentRepository) EndRent(ctx context.Context, rentID int64, lat, long float64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	tx, err := r.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// set timeend
	query := `UPDATE rents SET time_end = $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, query, time.Now(), rentID); err != nil {
		return err
	}

	// update transport
	var transportID int64
	query = `SELECT transport_id FROM rents WHERE id = $1`
	if err := r.QueryRowContext(ctx, query, rentID).Scan(&transportID); err != nil {
		return err
	}

	query = `UPDATE transports SET can_be_rented = TRUE, latitude = $1, longitude = $2 WHERE id = $3`
	if _, err := tx.ExecContext(ctx, query, lat, long, transportID); err != nil {
		return err
	}

	//set final price

	if err := tx.Commit(); err != nil {
		return err
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

func (r *RentRepository) Update(ctx context.Context, rent *entity.Rent) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE rents SET transport_id = $1, user_id = $2, time_start = $3, time_end = $4, price_of_unit = $5, price_type = $6, final_price = $7 WHERE id = $8`
	if _, err := r.ExecContext(ctx, query, rent.TransportID, rent.UserID, rent.TimeStart, rent.TimeEnd, rent.PriceOfUnit, rent.PriceType, rent.FinalPrice, rent.ID); err != nil {
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
