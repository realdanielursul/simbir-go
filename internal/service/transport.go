package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/entity"
	"github.com/realdanielursul/simbir-go/internal/repository"
)

type TransportService struct {
	transportRepo repository.Transport
}

func NewTransportService(transportRepo repository.Transport) *TransportService {
	return &TransportService{
		transportRepo: transportRepo,
	}
}

func (s *TransportService) CreateTransport(ctx context.Context, userID int64, input *TransportInput) (int64, error) {
	// validate data

	// check identifier uniqueness
	transport, err := s.transportRepo.GetByIdentifier(ctx, input.Identifier)
	if err != nil {
		return -1, err
	}

	if transport != nil {
		return -1, ErrIdentifierAlreadyExists
	}

	id, err := s.transportRepo.Create(ctx, &entity.Transport{
		OwnerID:       userID,
		CanBeRented:   input.CanBeRented,
		TransportType: input.TransportType,
		Model:         input.Model,
		Color:         input.Color,
		Identifier:    input.Identifier,
		Description:   input.Description,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		MinutePrice:   int64(input.MinutePrice * 100),
		DayPrice:      int64(input.DayPrice * 100),
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *TransportService) GetTransport(ctx context.Context, id int64) (*TransportOutput, error) {
	transport, err := s.transportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if transport == nil {
		return nil, ErrTransportNotFound
	}

	return &TransportOutput{
		ID:            transport.ID,
		OwnerID:       transport.OwnerID,
		CanBeRented:   transport.CanBeRented,
		TransportType: transport.TransportType,
		Model:         transport.Model,
		Color:         transport.Color,
		Identifier:    transport.Identifier,
		Description:   transport.Description,
		Latitude:      transport.Latitude,
		Longitude:     transport.Longitude,
		MinutePrice:   float64(transport.MinutePrice) / 100,
		DayPrice:      float64(transport.DayPrice) / 100,
		CreatedAt:     transport.CreatedAt,
		UpdatedAt:     transport.UpdatedAt,
	}, nil
}

func (s *TransportService) ListTransport(ctx context.Context, transportType string, count, start int) ([]TransportOutput, error) {
	transports, err := s.transportRepo.ListByType(ctx, transportType, count, start)
	if err != nil {
		return nil, err
	}

	transportsOutput := make([]TransportOutput, 0, len(transports))
	for _, transport := range transports {
		transportOutput := TransportOutput{
			ID:            transport.ID,
			OwnerID:       transport.OwnerID,
			CanBeRented:   transport.CanBeRented,
			TransportType: transport.TransportType,
			Model:         transport.Model,
			Color:         transport.Color,
			Identifier:    transport.Identifier,
			Description:   transport.Description,
			Latitude:      transport.Latitude,
			Longitude:     transport.Longitude,
			MinutePrice:   float64(transport.MinutePrice) / 100,
			DayPrice:      float64(transport.DayPrice) / 100,
			CreatedAt:     transport.CreatedAt,
			UpdatedAt:     transport.UpdatedAt,
		}

		transportsOutput = append(transportsOutput, transportOutput)
	}

	return transportsOutput, nil
}

func (s *TransportService) ListTransportByOwner(ctx context.Context, ownerID int64, count, start int) ([]TransportOutput, error) {
	transports, err := s.transportRepo.ListByOwner(ctx, ownerID, count, start)
	if err != nil {
		return nil, err
	}

	transportsOutput := make([]TransportOutput, 0, len(transports))
	for _, transport := range transports {
		transportOutput := TransportOutput{
			ID:            transport.ID,
			OwnerID:       transport.OwnerID,
			CanBeRented:   transport.CanBeRented,
			TransportType: transport.TransportType,
			Model:         transport.Model,
			Color:         transport.Color,
			Identifier:    transport.Identifier,
			Description:   transport.Description,
			Latitude:      transport.Latitude,
			Longitude:     transport.Longitude,
			MinutePrice:   float64(transport.MinutePrice) / 100,
			DayPrice:      float64(transport.DayPrice) / 100,
			CreatedAt:     transport.CreatedAt,
			UpdatedAt:     transport.UpdatedAt,
		}

		transportsOutput = append(transportsOutput, transportOutput)
	}

	return transportsOutput, nil
}

// func (s *TransportService) ListTransportAvailable(ctx context.Context, lat, long, radius float64, transportType string) ([]TransportOutput, error)

func (s *TransportService) UpdateTransport(ctx context.Context, userID, id int64, input *TransportInput) error {
	// validate data

	// check identifier uniqueness
	transport, err := s.transportRepo.GetByIdentifier(ctx, input.Identifier)
	if err != nil {
		return err
	}

	if transport != nil {
		return ErrIdentifierAlreadyExists
	}

	// check if user is owner
	if userID != transport.OwnerID {
		return ErrAccessDenied
	}

	err = s.transportRepo.Update(ctx, &entity.Transport{
		ID:          id,
		OwnerID:     userID,
		CanBeRented: input.CanBeRented,
		Model:       input.Model,
		Color:       input.Color,
		Identifier:  input.Identifier,
		Description: input.Description,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		MinutePrice: int64(input.MinutePrice * 100),
		DayPrice:    int64(input.DayPrice * 100),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *TransportService) DeleteTransport(ctx context.Context, userID, id int64) error {
	transport, err := s.transportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transport == nil {
		return ErrTransportNotFound
	}

	// check if user is owner
	if userID != transport.OwnerID {
		return ErrAccessDenied
	}

	err = s.transportRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
