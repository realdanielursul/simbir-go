package service

import (
	"context"

	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/realdanielursul/simbir-go/pkg/hasher"
)

type AdminAccountService struct {
	accountRepo    repository.Account
	tokenRepo      repository.Token
	passwordHasher hasher.PasswordHasher
}

func NewAdminAccountService(accountRepo repository.Account, tokenRepo repository.Token, passwordHasher hasher.PasswordHasher) *AdminAccountService {
	return &AdminAccountService{
		accountRepo:    accountRepo,
		tokenRepo:      tokenRepo,
		passwordHasher: passwordHasher,
	}
}

func (s *AdminAccountService) CreateAccount(ctx context.Context, input *AdminAccountInput) (int64, error) {
	panic("")
}

func (s *AdminAccountService) GetAccount(ctx context.Context, id int64) (*AdminAccountOutput, error) {
	panic("")
}

func (s *AdminAccountService) ListAccounts(ctx context.Context, count, start int) ([]AdminAccountOutput, error) {
	panic("")
}

func (s *AdminAccountService) UpdateAccount(ctx context.Context, id int64, input *AdminAccountInput) error {
	panic("")
}

func (s *AdminAccountService) DeleteAccount(ctx context.Context, id int64) error {
	panic("")
}
