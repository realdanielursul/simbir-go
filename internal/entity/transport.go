package entity

import "time"

type Transport struct {
	ID            int64     `db:"id"`
	OwnerID       int64     `db:"owner_id"`
	CanBeRented   bool      `db:"can_be_rented"`
	TransportType string    `db:"transport_type"`
	Model         string    `db:"model"`
	Color         string    `db:"color"`
	Identifier    string    `db:"identifier"`
	Description   *string   `db:"description"`
	Latitude      float64   `db:"latitude"`
	Longitude     float64   `db:"longitude"`
	MinutePrice   int64     `db:"minute_price"`
	DayPrice      int64     `db:"day_price"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
