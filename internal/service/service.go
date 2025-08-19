package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/realdanielursul/simbir-go/pkg/hasher"
)

type AccountInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountOutput struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Account interface {
	SignUp(ctx context.Context, input *AccountInput) (int64, error)
	SignIn(ctx context.Context, input *AccountInput) (string, error)
	SignOut(ctx context.Context, tokenString string) error
	GetAccount(ctx context.Context, id int64) (*AccountOutput, error)
	UpdateAccount(ctx context.Context, id int64, input *AccountInput) error
	ValidateToken(ctx context.Context, tokenString string) (int64, bool, error)
}

type AdminAccountInput struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	IsAdmin  bool    `json:"isAdmin"`
	Balance  float64 `json:"balance"`
}

type AdminAccountOutput struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"isAdmin"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AdminAccount interface {
	CreateAccount(ctx context.Context, input *AdminAccountInput) (int64, error)
	GetAccount(ctx context.Context, id int64) (*AdminAccountOutput, error)
	ListAccounts(ctx context.Context, count, start int) ([]AdminAccountOutput, error)
	UpdateAccount(ctx context.Context, id int64, input *AdminAccountInput) error
	DeleteAccount(ctx context.Context, id int64) error
}

type TransportInput struct {
	CanBeRented   bool    `json:"canBeRented"`
	TransportType string  `json:"transportType"`
	Model         string  `json:"model"`
	Color         string  `json:"color"`
	Identifier    string  `json:"identifier"`
	Description   *string `json:"description"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	MinutePrice   float64 `json:"minutePrice"`
	DayPrice      float64 `json:"dayPrice"`
}

type TransportOutput struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"ownerId"`
	CanBeRented   bool      `json:"canBeRented"`
	TransportType string    `json:"transportType"`
	Model         string    `json:"model"`
	Color         string    `json:"color"`
	Identifier    string    `json:"identifier"`
	Description   *string   `json:"description,omitempty"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	MinutePrice   float64   `json:"minutePrice"`
	DayPrice      float64   `json:"dayPrice"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Transport interface {
	CreateTransport(ctx context.Context, userID int64, input *TransportInput) (int64, error)
	GetTransport(ctx context.Context, id int64) (*TransportOutput, error)
	ListTransport(ctx context.Context, transportType string, count, start int) ([]TransportOutput, error)
	ListTransportByOwner(ctx context.Context, ownerID int64, count, start int) ([]TransportOutput, error)
	ListTransportByAvailability(ctx context.Context, lat, long, radius float64, transportType string) ([]TransportOutput, error)
	UpdateTransport(ctx context.Context, userID, id int64, input *TransportInput) error
	DeleteTransport(ctx context.Context, userID, id int64) error
}

type AdminTransportInput struct {
	OwnerID       int64   `json:"ownerId"`
	CanBeRented   bool    `json:"canBeRented"`
	TransportType string  `json:"transportType"`
	Model         string  `json:"model"`
	Color         string  `json:"color"`
	Identifier    string  `json:"identifier"`
	Description   *string `json:"description,omitempty"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	MinutePrice   float64 `json:"minutePrice"`
	DayPrice      float64 `json:"dayPrice"`
}

type AdminTransport interface {
	CreateTransport(ctx context.Context, input *AdminTransportInput) (int64, error)
	UpdateTransport(ctx context.Context, id int64, input *AdminTransportInput) error
	DeleteTransport(ctx context.Context, id int64) error
}

type Payment interface {
	UpdateBalance(ctx context.Context, accountID int64, amount float64) error
	BillingWorker(ctx context.Context)
	ProcessBilling(ctx context.Context)
}

type RentOutput struct {
	ID           int64      `json:"id"`
	TransportID  int64      `json:"transportId"`
	UserID       int64      `json:"userId"`
	TimeStart    time.Time  `json:"timeStart"`
	TimeEnd      *time.Time `json:"timeEnd,omitempty"`
	PriceOfUnit  int64      `json:"priceOfUnit"`
	PriceType    string     `json:"priceType"`
	FinalPrice   *int64     `json:"finalPrice,omitempty"`
	LastBilledAt time.Time  `json:"lastBilledAt"`
}

type Rent interface {
	StartRent(ctx context.Context, userID, transportID int64, rentType string) (int64, error)
	StopRent(ctx context.Context, userID, id int64, lat, long float64) error
	GetRent(ctx context.Context, id int64) (*RentOutput, error)
	ListRentsByAccount(ctx context.Context, accountID int64) ([]RentOutput, error)
	ListRentsByTransport(ctx context.Context, transportID int64) ([]RentOutput, error)
}

type ServicesDependencies struct {
	Repos    *repository.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

// type Services struct {
// 	Account        Account
// 	AdminAccount   AdminAccount
// 	Transport      Transport
// 	AdminTransport AdminTransport
// }

// func NewServices(deps ServicesDependencies) *Services {
// 	return &Services{
// 		Account:        NewAccountService(deps.Repos.Account, deps.Repos.Token, deps.Hasher, deps.SignKey, deps.TokenTTL),
// 		AdminAccount:   NewAdminAccountService(deps.Repos.Account, deps.Repos.Token, deps.Hasher),
// 		Transport:      NewTransportService(deps.Repos.Transport),
// 		AdminTransport: NewAdminTransportService(deps.Repos.Transport),
// 	}
// }
