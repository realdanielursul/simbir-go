package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

const operationTimeout = 3 * time.Second

type Account interface {
	Create(ctx context.Context, account *entity.Account) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
	GetByUsernameAndPassword(ctx context.Context, username, password string) (*entity.Account, error)
	List(ctx context.Context, count, start int) ([]entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id int64) error
}

type Token interface {
	Create(ctx context.Context, token *entity.Token) error
	Get(ctx context.Context, tokenString string) (*entity.Token, error)
	Invalidate(ctx context.Context, tokenString string) error
	InvalidateAll(ctx context.Context, id int64) error
}

type Transport interface {
	Create(ctx context.Context, transport *entity.Transport) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Transport, error)
	GetByIdentifier(ctx context.Context, identifier string) (*entity.Transport, error)
	ListByType(ctx context.Context, transportType string, count, start int) ([]entity.Transport, error)
	ListByOwner(ctx context.Context, ownerID int64, count, start int) ([]entity.Transport, error)
	ListByAvailability(ctx context.Context, lat, long, radius float64, transportType string) ([]entity.Transport, error)
	Update(ctx context.Context, transport *entity.Transport) error
	ChangeAvailability(ctx context.Context, id int64, can_be_rented bool) error
	Delete(ctx context.Context, id int64) error
}

type Payment interface {
	UpdateBalance(ctx context.Context, accountID, amount int64) error
}

type Rent interface {
	StartRent(ctx context.Context, rent *entity.Rent) (int64, error)
	EndRent(ctx context.Context, id int64, lat, long float64) error
	GetByID(ctx context.Context, id int64) (*entity.Rent, error)
	GetHistoryByUser(ctx context.Context, userID int64) ([]entity.Rent, error)
	GetHistoryByTransport(ctx context.Context, transportID int64) ([]entity.Rent, error)
	ListActive(ctx context.Context) ([]entity.Rent, error)
	Update(ctx context.Context, rent *entity.Rent) error
	UpdateLastBilledTime(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
}

type Repositories struct {
	Account
	Token
	Transport
	Payment
	Rent
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Account:   NewAccountRepository(db),
		Token:     NewTokenRepository(db),
		Transport: NewTransportRepository(db),
		Payment:   NewPaymentRepository(db),
		Rent:      NewRentRepository(db),
	}
}
