package feedback

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines behavior for submitting ratings/feedback.
type Service interface {
	SetRating(ctx context.Context, orderID string, req model.RatingRequest) (model.RatingResponse, error)
}


