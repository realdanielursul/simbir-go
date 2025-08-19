package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/repository"
	"github.com/sirupsen/logrus"
)

type PaymentService struct {
	accountRepo   repository.Account
	transportRepo repository.Transport
	paymentRepo   repository.Payment
	rentRepo      repository.Rent
}

func NewPaymentService(accountRepo repository.Account, paymentRepo repository.Payment, transportRepo repository.Transport, rentRepo repository.Rent) *PaymentService {
	return &PaymentService{
		accountRepo:   accountRepo,
		transportRepo: transportRepo,
		paymentRepo:   paymentRepo,
		rentRepo:      rentRepo,
	}
}

func (s *PaymentService) UpdateBalance(ctx context.Context, accountID int64, amount float64) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account == nil {
		return ErrAccountNotFound
	}

	if err := s.paymentRepo.UpdateBalance(ctx, accountID, int64(amount*100)); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) BillingWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.ProcessBilling(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *PaymentService) ProcessBilling(ctx context.Context) {
	rents, err := s.rentRepo.ListActive(ctx)
	if err != nil {
		logrus.Errorf("billing error: %w", err)
	}

	for _, rent := range rents {
		account, err := s.accountRepo.GetByID(ctx, rent.UserID)
		if err != nil {
			logrus.Errorf("billing error: %w", err)
		}

		elapsed := time.Now().UTC().Sub(rent.LastBilledAt)

		switch rent.PriceType {
		case "Minutes":
			if account.Balance < rent.PriceOfUnit {
				if err := s.rentRepo.EndRent(ctx, rent.ID, 0, 0); err != nil {
					logrus.Errorf("billing error: %w", err)
				}

				if err := s.transportRepo.ChangeAvailability(ctx, rent.TransportID, true); err != nil {
					logrus.Errorf("billing error: %w", err)
				}
			}

			if int64(elapsed.Minutes()) == 1 {
				if err := s.paymentRepo.UpdateBalance(ctx, rent.UserID, rent.PriceOfUnit); err != nil {
					logrus.Errorf("billing error: %w", err)
				}

				if err := s.rentRepo.UpdateLastBilledTime(ctx, rent.ID); err != nil {
					logrus.Errorf("billing error: %w", err)
				}
			}
		case "Days":
			if account.Balance < rent.PriceOfUnit {
				if err := s.rentRepo.EndRent(ctx, rent.ID, 0, 0); err != nil {
					logrus.Errorf("billing error: %w", err)
				}

				if err := s.transportRepo.ChangeAvailability(ctx, rent.TransportID, true); err != nil {
					logrus.Errorf("billing error: %w", err)
				}
			}

			if int64(elapsed.Hours()/24) == 1 {
				if err := s.paymentRepo.UpdateBalance(ctx, rent.UserID, rent.PriceOfUnit); err != nil {
					logrus.Errorf("billing error: %w", err)
				}

				if err := s.rentRepo.UpdateLastBilledTime(ctx, rent.ID); err != nil {
					logrus.Errorf("billing error: %w", err)
				}
			}
		}
	}
}
