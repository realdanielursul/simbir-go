package repository

import (
	"context"
	"database/sql"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

type TransportRepository struct {
	*sqlx.DB
}

func NewTransportRepository(db *sqlx.DB) *TransportRepository {
	return &TransportRepository{db}
}

func (r *TransportRepository) Create(ctx context.Context, transport *entity.Transport) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var id int64
	query := `INSERT INTO transports (owner_id, can_be_rented, transport_type, model, color, identifier, description, latitude, longitude, minute_price, day_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	if err := r.QueryRowContext(ctx, query, transport.OwnerID, transport.CanBeRented, transport.TransportType, transport.Model, transport.Color, transport.Identifier, transport.Description, transport.Latitude, transport.Longitude, transport.MinutePrice, transport.DayPrice).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *TransportRepository) GetByID(ctx context.Context, id int64) (*entity.Transport, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var transport entity.Transport
	query := `SELECT * FROM transports WHERE id = $1`
	if err := r.QueryRowxContext(ctx, query, id).StructScan(&transport); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &transport, nil
}

func (r *TransportRepository) GetByIdentifier(ctx context.Context, identifier string) (*entity.Transport, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var transport entity.Transport
	query := `SELECT * FROM transports WHERE identifier = $1`
	if err := r.QueryRowxContext(ctx, query, identifier).StructScan(&transport); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &transport, nil
}

func (r *TransportRepository) ListByType(ctx context.Context, transportType string, count, start int) ([]entity.Transport, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	transports := make([]entity.Transport, 0, count)

	var query string
	var rows *sqlx.Rows
	var err error

	if transportType == "All" {
		query = `SELECT * FROM transports WHERE can_be_rented = TRUE ORDER BY id`
		rows, err = r.QueryxContext(ctx, query)
	} else {
		query = `SELECT * FROM transports WHERE can_be_rented = TRUE AND transport_type = $1 ORDER BY id`
		rows, err = r.QueryxContext(ctx, query, transportType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transport entity.Transport
		if err := rows.StructScan(&transport); err != nil {
			return nil, err
		}

		transports = append(transports, transport)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transports, nil
}

func (r *TransportRepository) ListByOwner(ctx context.Context, ownerID int64, count, start int) ([]entity.Transport, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	transports := make([]entity.Transport, 0, count)
	query := `SELECT * FROM transports WHERE owner_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.QueryxContext(ctx, query, ownerID, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transport entity.Transport
		if err := rows.StructScan(&transport); err != nil {
			return nil, err
		}

		transports = append(transports, transport)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transports, nil
}

func (r *TransportRepository) ListByAvailability(ctx context.Context, lat, long, radius float64, transportType string) ([]entity.Transport, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	transports := make([]entity.Transport, 0, 100)

	var query string
	var rows *sqlx.Rows
	var err error

	if transportType == "All" {
		query = `SELECT * FROM transports WHERE can_be_rented = TRUE ORDER BY id`
		rows, err = r.QueryxContext(ctx, query)
	} else {
		query = `SELECT * FROM transports WHERE can_be_rented = TRUE AND transport_type = $1 ORDER BY id`
		rows, err = r.QueryxContext(ctx, query, transportType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transport entity.Transport
		if err := rows.StructScan(&transport); err != nil {
			return nil, err
		}

		if isAvailable(lat, long, transport.Latitude, transport.Longitude, radius) {
			transports = append(transports, transport)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transports, nil
}

func (r *TransportRepository) Update(ctx context.Context, transport *entity.Transport) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `UPDATE transports SET owner_id = $1, can_be_rented = $2, transport_type = $3, model = $4, color = $5, identifier = $6, description = $7, latitude = $8, longitude = $9, minute_price = $10, day_price = $11, updated_at = $12 WHERE id = $13`
	if _, err := r.ExecContext(ctx, query, transport.OwnerID, transport.CanBeRented, transport.TransportType, transport.Model, transport.Color, transport.Identifier, transport.Description, transport.Latitude, transport.Longitude, transport.MinutePrice, transport.DayPrice, transport.UpdatedAt, transport.ID); err != nil {
		return err
	}

	return nil
}

func (r *TransportRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	query := `DELETE FROM transports WHERE id = $1`
	if _, err := r.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func isAvailable(lat1, long1, lat2, long2, radius float64) bool {
	x := lat1 - lat2
	y := long1 - long2
	r := math.Sqrt(x*x + y*y)

	return r < radius
}
