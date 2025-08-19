package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/entity"
	"github.com/realdanielursul/simbir-go/internal/repository"
)

type RentService struct {
	accountRepo   repository.Account
	transportRepo repository.Transport
	rentRepo      repository.Rent
}

func NewRentService(accountRepo repository.Account, transportRepo repository.Transport, rentRepo repository.Rent) *RentService {
	return &RentService{
		accountRepo:   accountRepo,
		transportRepo: transportRepo,
		rentRepo:      rentRepo,
	}
}

func (s *RentService) StartRent(ctx context.Context, userID, transportID int64, rentType string) (int64, error) {
	// check balance is good
	//get acc
	account, err := s.accountRepo.GetByID(ctx, userID)
	if err != nil {
		return -1, err
	}

	if account == nil {
		return -1, ErrAccountNotFound
	}

	//get transport
	transport, err := s.transportRepo.GetByID(ctx, transportID)
	if err != nil {
		return -1, err
	}

	if transport == nil {
		return -1, ErrTransportNotFound
	}

	if rentType == "Minutes" && account.Balance < transport.MinutePrice {
		return -1, ErrNotEnoughMoney
	} else if rentType == "Days" && account.Balance < transport.DayPrice {
		return -1, ErrNotEnoughMoney
	} else {
		return -1, ErrInvalidRentType
	}

	var priceOfUnit int64
	if rentType == "Minutes" {
		priceOfUnit = transport.MinutePrice
	} else {
		priceOfUnit = transport.DayPrice
	}

	id, err := s.rentRepo.StartRent(ctx, &entity.Rent{
		TransportID: transportID,
		UserID:      userID,
		TimeStart:   time.Now().UTC(),
		TimeEnd:     nil,
		PriceOfUnit: priceOfUnit,
		PriceType:   rentType,
		FinalPrice:  nil,
	})

	_, err = s.rentRepo.GetByID(ctx, id)
	if err != nil {
		return -1, err
	}

	return -1, nil
}

func (s *RentService) StopRent(ctx context.Context, userID, id int64, lat, long float64) error {
	panic("")
}
func (s *RentService) GetRent(ctx context.Context, id int64) (*RentOutput, error)
func (s *RentService) ListRentsByAccount(ctx context.Context, accountID int64) ([]RentOutput, error)
func (s *RentService) ListRentsByTransport(ctx context.Context, transportID int64) ([]RentOutput, error)
