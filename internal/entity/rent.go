package entity

import "time"

type Rent struct {
	ID          int64      `db:"id"`
	TransportID int64      `db:"transport_id"`
	UserID      int64      `db:"user_id"`
	TimeStart   time.Time  `db:"time_start"`
	TimeEnd     *time.Time `db:"time_end"`
	PriceOfUnit int64      `db:"price_of_unit"`
	PriceType   string     `db:"price_type"`
	FinalPrice  *int64     `db:"final_price"`
}
