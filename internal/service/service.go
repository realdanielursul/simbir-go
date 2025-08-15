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
	SignIn(ctx context.Context, id int64) (string, error)
	SignOut(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (*AccountOutput, error)
	UpdateAccount(ctx context.Context, id int64, input *AccountInput) error
	ValidateToken(ctx context.Context, tokenString string) (int, bool, error)
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

type ServicesDependencies struct {
	Repos    *repository.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Account      Account
	AdminAccount AdminAccount
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Account:      NewAccountService(deps.Repos.Account, deps.Repos.Token, deps.Hasher, deps.SignKey, deps.TokenTTL),
		AdminAccount: NewAdminAccountService(deps.Repos.Account, deps.Repos.Token, deps.Hasher),
	}
}
