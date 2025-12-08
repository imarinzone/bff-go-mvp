package support

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data based on swagger example.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) GetSupport(ctx context.Context, orderID string) (model.SupportResponse, error) {
	_ = ctx
	_ = orderID

	return model.SupportResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "CANCELLED",
			Mode:   "RESERVATION",
		},
		Name:  "BlueCharge Support Team",
		Phone: "18001080",
		Email: "support@bluechargenet-aggregator.io",
		URL:   "https://support.bluechargenet-aggregator.io/ticket/SUP-20250730-001",
		Hours: "Mon–Sun 24/7 IST",
		Channels: []string{
			"phone",
			"email",
			"web",
			"chat",
		},
	}, nil
}
