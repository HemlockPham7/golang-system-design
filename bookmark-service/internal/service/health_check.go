package service

import (
	"context"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/repository"
)

type HealthCheck interface {
	HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error)
}

type healthCheck struct {
	serviceName string
	instanceId  string
	pingRepo    repository.Ping
}

func NewHealthCheck(serviceName, instanceId string, pingRepo repository.Ping) HealthCheck {
	return &healthCheck{
		serviceName: serviceName,
		instanceId:  instanceId,
		pingRepo:    pingRepo,
	}
}

func (s *healthCheck) HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error) {
	if err := s.pingRepo.CheckHealth(ctx); err != nil {
		return nil, err
	}
	return &model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  s.instanceId,
	}, nil
}
