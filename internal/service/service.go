package service

import (
	"context"

	grpcclient "github.com/KirillKotovsky/location_proto/proto/pkg/grpcclient"
	"github.com/KirillKotovsky/location_user_api/internal/model"
	"github.com/KirillKotovsky/location_user_api/internal/repository"
)

type Service struct {
	repo       repository.Repository
	grpcClient grpcclient.Client
}

func New(repo repository.Repository, grpc grpcclient.Client) *Service {
	return &Service{repo: repo, grpcClient: grpc}
}

func (s *Service) UpdateLocation(ctx context.Context, req model.LocationUpdateRequest) error {
	location := req.ToEntity()

	// Сохраняем текущую локацию
	if err := s.repo.SaveLocation(ctx, location); err != nil {
		return err
	}

	// Отправляем в gRPC-сервис
	if err := s.grpcClient.StoreLocation(ctx, location.ToProto()); err != nil {
		return err
	}

	return nil
}

func (s *Service) FindNearbyUsers(ctx context.Context, req model.NearbyUsersRequest) ([]model.UserLocation, error) {
	return s.repo.FindNearby(ctx, req)
}
