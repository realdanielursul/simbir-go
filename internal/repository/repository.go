package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/realdanielursul/simbir-go/internal/entity"
)

// AccountRepository отвечает за CRUD аккаунтов и аутентификацию
type Account interface {
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
	Create(ctx context.Context, account *entity.Account) (int64, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, start, count int) ([]entity.Account, error)
	Authenticate(ctx context.Context, username, password string) (*entity.Account, error)
	UpdateBalance(ctx context.Context, id int64, amount float64) error
}

type Token interface {
	CreateToken(ctx context.Context, token *entity.Token) error
	GetToken(ctx context.Context, tokenString string) (*entity.Token, error)
	InvalidateUserTokens(ctx context.Context, login string) error
}

type Transport interface {
	GetByID(ctx context.Context, id int64) (*entity.Transport, error)
	Create(ctx context.Context, transport *entity.Transport) (int64, error)
	Update(ctx context.Context, transport *entity.Transport) error
	Delete(ctx context.Context, id int64) error

	ListAll(ctx context.Context, start, count int, transportType string) ([]entity.Transport, error)
	SearchAvailable(ctx context.Context, lat, long, radius float64, transportType string) ([]entity.Transport, error)
	ListByOwner(ctx context.Context, ownerID int64, start, count int) ([]entity.Transport, error)
}

type Rent interface {
	GetByID(ctx context.Context, id int64) (*entity.Rent, error)
	Create(ctx context.Context, rent *entity.Rent) (int64, error)
	Update(ctx context.Context, rent *entity.Rent) error
	Delete(ctx context.Context, id int64) error

	GetHistoryByUser(ctx context.Context, userID int64) ([]entity.Rent, error)
	GetHistoryByTransport(ctx context.Context, transportID int64) ([]entity.Rent, error)
	GetActiveByUser(ctx context.Context, userID int64) ([]entity.Rent, error)
	GetAvailableForRent(ctx context.Context, lat, long, radius float64, transportType string) ([]entity.Transport, error)
	EndRent(ctx context.Context, rentID int64, lat, long float64) error
}

type Payment interface {
	AddBalance(ctx context.Context, accountID int64, amount float64) error
}

type Repositories struct {
	Account
	Token
	Transport
	Rent
	Payment
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		// initialize repos
	}
}
