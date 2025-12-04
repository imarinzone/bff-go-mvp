package estimate

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines the behavior for the estimate use case.
type Service interface {
	Estimate(ctx context.Context, req model.EstimateRequest) (model.EstimateResponse, error)
}


