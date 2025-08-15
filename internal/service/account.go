package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/realdanielursul/simbir-go/pkg/hasher"
)

type AccountService struct {
	accountRepo    repository.Account
	tokenRepo      repository.Token
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAccountService(accountRepo repository.Account, tokenRepo repository.Token, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AccountService {
	return &AccountService{
		accountRepo:    accountRepo,
		tokenRepo:      tokenRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AccountService) SignUp(ctx context.Context, input *AccountInput) (int64, error) {
	panic("")
}

func (s *AccountService) SignIn(ctx context.Context, id int64) (string, error) {
	panic("")
}

func (s *AccountService) SignOut(ctx context.Context, id int64) error {
	panic("")
}

func (s *AccountService) GetAccount(ctx context.Context, id int64) (*AccountOutput, error) {
	panic("")
}

func (s *AccountService) UpdateAccount(ctx context.Context, id int64, input *AccountInput) error {
	panic("")
}

func (s *AccountService) ValidateToken(ctx context.Context, tokenString string) (int, bool, error) {
	panic("")
}
