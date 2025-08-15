package service

import (
	"context"

	"github.com/realdanielursul/simbir-go/internal/entity"
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
	// check username uniqueness
	account, err := s.accountRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return -1, err
	}

	if account != nil {
		return -1, ErrInvalidUsername
	}

	// create account in db
	id, err := s.accountRepo.Create(ctx, &entity.Account{
		Username:     input.Username,
		PasswordHash: s.passwordHasher.Hash(input.Password),
		IsAdmin:      input.IsAdmin,
		Balance:      int64(input.Balance * 100),
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *AdminAccountService) GetAccount(ctx context.Context, id int64) (*AdminAccountOutput, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return &AdminAccountOutput{
		ID:        account.ID,
		Username:  account.Username,
		IsAdmin:   account.IsAdmin,
		Balance:   float64(account.Balance) / 100,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}

func (s *AdminAccountService) ListAccounts(ctx context.Context, count, start int) ([]AdminAccountOutput, error) {
	accounts, err := s.accountRepo.List(ctx, count, start)
	if err != nil {
		return nil, err
	}

	accountOutputs := make([]AdminAccountOutput, 0, 100)
	for _, account := range accounts {
		accountOutput := AdminAccountOutput{
			ID:        account.ID,
			Username:  account.Username,
			IsAdmin:   account.IsAdmin,
			Balance:   float64(account.Balance) / 100,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		}

		accountOutputs = append(accountOutputs, accountOutput)
	}

	return accountOutputs, nil
}

func (s *AdminAccountService) UpdateAccount(ctx context.Context, id int64, input *AdminAccountInput) error {
	// check username uniqueness
	account, err := s.accountRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return err
	}

	if account != nil {
		return ErrInvalidUsername
	}

	err = s.accountRepo.Update(ctx, &entity.Account{
		ID:           id,
		Username:     input.Username,
		PasswordHash: s.passwordHasher.Hash(input.Password),
		IsAdmin:      input.IsAdmin,
		Balance:      int64(input.Balance * 100),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminAccountService) DeleteAccount(ctx context.Context, id int64) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if account == nil {
		return ErrAccountNotFound
	}

	if err := s.accountRepo.Delete(ctx, id); err != nil {
		return err
	}

	return err
}
