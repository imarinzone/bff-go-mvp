package estimate

import (
	"context"

	"go.uber.org/zap"

	"bff-go-mvp/internal/model"
)

// FallbackService wraps a primary service and falls back to a secondary service on errors
type FallbackService struct {
	primary  Service
	fallback Service
	logger   *zap.Logger
}

// NewFallbackService creates a new fallback service
func NewFallbackService(primary, fallback Service, logger *zap.Logger) *FallbackService {
	return &FallbackService{
		primary:  primary,
		fallback: fallback,
		logger:   logger,
	}
}

// Estimate implements the Service interface with fallback logic
func (s *FallbackService) Estimate(ctx context.Context, req model.EstimateRequest) (model.EstimateResponse, error) {
	resp, err := s.primary.Estimate(ctx, req)
	if err != nil {
		s.logger.Warn("Primary service failed, falling back to secondary service", zap.Error(err))
		return s.fallback.Estimate(ctx, req)
	}
	return resp, nil
}

