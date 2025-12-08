package feedback

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data based on swagger example.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) SetRating(ctx context.Context, orderID string, req model.RatingRequest) (model.RatingResponse, error) {
	_ = ctx
	_ = orderID
	_ = req

	return model.RatingResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "COMPLETED",
			Mode:   "RESERVATION",
		},
		FeedbackForm: &model.FeedbackForm{
			URL:          "https://example-bpp.com/feedback/portal",
			MIMEType:     "application/xml",
			SubmissionID: "feedback-123e4567-e89b-12d3-a456-426614174000",
		},
	}, nil
}
