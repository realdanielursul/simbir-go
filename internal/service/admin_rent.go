package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/entity"
	"github.com/realdanielursul/simbir-go/internal/repository"
)

type AdminRentService struct {
	accountRepo   repository.Account
	paymentRepo   repository.Payment
	transportRepo repository.Transport
	rentRepo      repository.Rent
}

func NewAdminRentService(accountRepo repository.Account, paymentRepo repository.Payment, transportRepo repository.Transport, rentRepo repository.Rent) *AdminRentService {
	return &AdminRentService{
		accountRepo:   accountRepo,
		paymentRepo:   paymentRepo,
		transportRepo: transportRepo,
		rentRepo:      rentRepo,
	}
}

func (s *AdminRentService) StartRent(ctx context.Context, input *AdminRentInput) (int64, error) {
	// check balance is good
	//get acc
	account, err := s.accountRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return -1, err
	}

	if account == nil {
		return -1, ErrAccountNotFound
	}

	//get transport
	transport, err := s.transportRepo.GetByID(ctx, input.TransportID)
	if err != nil {
		return -1, err
	}

	if transport == nil {
		return -1, ErrTransportNotFound
	}

	var priceOfUnit int64
	if input.PriceType == "Minutes" {
		priceOfUnit = transport.MinutePrice
	} else {
		priceOfUnit = transport.DayPrice
	}

	id, err := s.rentRepo.StartRent(ctx, &entity.Rent{
		TransportID: input.TransportID,
		UserID:      input.UserID,
		TimeStart:   time.Now().UTC(),
		TimeEnd:     nil,
		PriceOfUnit: priceOfUnit,
		PriceType:   input.PriceType,
		FinalPrice:  nil,
	})

	if err := s.transportRepo.ChangeAvailability(ctx, input.TransportID, false); err != nil {
		return -1, err
	}

	if err := s.paymentRepo.UpdateBalance(ctx, input.UserID, -priceOfUnit); err != nil {
		return -1, ErrRentNotFound
	}

	return id, nil
}

func (s *AdminRentService) EndRent(ctx context.Context, id int64, lat, long float64) error {
	rent, err := s.rentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if rent == nil {
		return ErrRentNotFound
	}

	if err := s.rentRepo.EndRent(ctx, id, lat, long); err != nil {
		return err
	}

	if err := s.transportRepo.ChangeAvailability(ctx, rent.TransportID, true); err != nil {
		return err
	}

	return nil
}

func (s *AdminRentService) GetRent(ctx context.Context, id int64) (*RentOutput, error) {
	rent, err := s.rentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if rent == nil {
		return nil, ErrRentNotFound
	}

	return &RentOutput{
		ID:          rent.ID,
		TransportID: rent.TransportID,
		UserID:      rent.UserID,
		TimeStart:   rent.TimeStart,
		TimeEnd:     rent.TimeEnd,
		PriceOfUnit: rent.PriceOfUnit,
		PriceType:   rent.PriceType,
		FinalPrice:  rent.FinalPrice,
	}, nil
}

func (s *AdminRentService) ListRentsByUser(ctx context.Context, userID int64) ([]RentOutput, error) {
	account, err := s.accountRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	rents, err := s.rentRepo.GetHistoryByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	rentsOutput := make([]RentOutput, 0, len(rents))
	for _, rent := range rents {
		rentOutput := RentOutput{
			ID:          rent.ID,
			TransportID: rent.TransportID,
			UserID:      rent.UserID,
			TimeStart:   rent.TimeStart,
			TimeEnd:     rent.TimeEnd,
			PriceOfUnit: rent.PriceOfUnit,
			PriceType:   rent.PriceType,
			FinalPrice:  rent.FinalPrice,
		}

		rentsOutput = append(rentsOutput, rentOutput)
	}

	return rentsOutput, nil
}

func (s *AdminRentService) ListRentsByTransport(ctx context.Context, transportID int64) ([]RentOutput, error) {
	transport, err := s.transportRepo.GetByID(ctx, transportID)
	if err != nil {
		return nil, err
	}

	if transport == nil {
		return nil, ErrTransportNotFound
	}

	rents, err := s.rentRepo.GetHistoryByTransport(ctx, transportID)
	if err != nil {
		return nil, err
	}

	rentsOutput := make([]RentOutput, 0, len(rents))
	for _, rent := range rents {
		rentOutput := RentOutput{
			ID:          rent.ID,
			TransportID: rent.TransportID,
			UserID:      rent.UserID,
			TimeStart:   rent.TimeStart,
			TimeEnd:     rent.TimeEnd,
			PriceOfUnit: rent.PriceOfUnit,
			PriceType:   rent.PriceType,
			FinalPrice:  rent.FinalPrice,
		}

		rentsOutput = append(rentsOutput, rentOutput)
	}

	return rentsOutput, nil
}

func (s *AdminRentService) DeleteRent(ctx context.Context, id int64) error {
	rent, err := s.rentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if rent == nil {
		return ErrRentNotFound
	}

	if err := s.rentRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
