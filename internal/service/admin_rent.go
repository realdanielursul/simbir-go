package service

import (
	"context"
)

type AdminRentService struct {
}

func NewAdminRentService() *AdminRentService {
	return &AdminRentService{}
}

func (s *AdminRentService) GetRent(ctx context.Context, id int64) (*RentOutput, error)
func (s *AdminRentService) ListRentsByUser(ctx context.Context, userID int64) ([]RentOutput, error)
func (s *AdminRentService) ListRentsByTransport(ctx context.Context, transportID int64) ([]RentOutput, error)
func (s *AdminRentService) StartRent(ctx context.Context, input *AdminRentInput) (int64, error)
func (s *AdminRentService) EndRent(ctx context.Context, id int64, lat, long float64) error
func (s *AdminRentService) DeleteRent(ctx context.Context, id int64) error
