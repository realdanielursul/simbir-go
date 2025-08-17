package service

import (
	"context"
	"time"

	"github.com/realdanielursul/simbir-go/internal/entity"
	"github.com/realdanielursul/simbir-go/internal/repository"
)

type AdminTransportService struct {
	transportRepo repository.Transport
}

func NewAdminTransportService(transportRepo repository.Transport) *AdminTransportService {
	return &AdminTransportService{
		transportRepo: transportRepo,
	}
}

func (s *AdminTransportService) CreateTransport(ctx context.Context, input *AdminTransportInput) (int64, error) {
	// check identifier uniqueness
	transport, err := s.transportRepo.GetByIdentifier(ctx, input.Identifier)
	if err != nil {
		return -1, err
	}

	if transport != nil {
		return -1, ErrIdentifierAlreadyExists
	}

	id, err := s.transportRepo.Create(ctx, &entity.Transport{
		OwnerID:       input.OwnerID,
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

func (s *AdminTransportService) UpdateTransport(ctx context.Context, id int64, input *AdminTransportInput) error {
	// check identifier uniqueness
	transport, err := s.transportRepo.GetByIdentifier(ctx, input.Identifier)
	if err != nil {
		return err
	}

	if transport != nil {
		return ErrIdentifierAlreadyExists
	}

	err = s.transportRepo.Update(ctx, &entity.Transport{
		ID:            id,
		OwnerID:       input.OwnerID,
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
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminTransportService) DeleteTransport(ctx context.Context, id int64) error {
	transport, err := s.transportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transport == nil {
		return ErrTransportNotFound
	}

	if err := s.transportRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
