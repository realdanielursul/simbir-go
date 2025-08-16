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
	// validate transport

	// check identifier uniqueness
	transport, err := s.transportRepo.GetByIdentifier(ctx, input.Identifier)
	if err != nil {
		return -1, err
	}

	if transport != nil {
		return -1, err
	}

	// create transport
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
	transports := make([]entity.Transport, 0, 100)
	var err error

	if transportType == "All" {
		transports, err = s.transportRepo.ListAll(ctx, count, start)
	} else {
		transports, err = s.transportRepo.ListByType(ctx, transportType, count, start)
	}

	if err != nil {
		return nil, err
	}

	transportOutputs := make([]TransportOutput, 0, 100)
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

		transportOutputs = append(transportOutputs, transportOutput)
	}

	return transportOutputs, nil
}

func (s *TransportService) ListTransportByOwner(ctx context.Context, ownerID int64, count, start int) ([]TransportOutput, error) {
	transports, err := s.transportRepo.ListByOwner(ctx, ownerID, count, start)
	if err != nil {
		return nil, err
	}

	transportOutputs := make([]TransportOutput, 0, 100)
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

		transportOutputs = append(transportOutputs, transportOutput)
	}

	return transportOutputs, nil
}

func (s *TransportService) ListTransportAvailable(ctx context.Context, lat, long, radius float64, transportType string) ([]TransportOutput, error) {
	transports, err := s.transportRepo.ListAvailable(ctx, lat, long, radius, transportType)
	if err != nil {
		return nil, err
	}

	transportOutputs := make([]TransportOutput, 0, 100)
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

		transportOutputs = append(transportOutputs, transportOutput)
	}

	return transportOutputs, nil
}

func (s *TransportService) UpdateTransport(ctx context.Context, userID, id int64, input *TransportInput) error {
	transport, err := s.transportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transport == nil {
		return err
	}

	if userID != transport.OwnerID {
		return ErrAccessDenied
	}

	err = s.transportRepo.Update(ctx, &entity.Transport{
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
		UpdatedAt:     time.Now(),
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

	if userID == transport.OwnerID {
		return ErrAccessDenied
	}

	err = s.transportRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
