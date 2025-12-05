package search

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

// Search implements the Service interface with fallback logic
func (s *FallbackService) Search(ctx context.Context, page, perPage int, req model.SearchRequest) (model.SearchResponse, error) {
	resp, err := s.primary.Search(ctx, page, perPage, req)
	if err != nil {
		s.logger.Warn("Primary service failed, falling back to secondary service", zap.Error(err))
		return s.fallback.Search(ctx, page, perPage, req)
	}
	return resp, nil
}
