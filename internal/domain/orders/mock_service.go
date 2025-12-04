package orders

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data that matches
// the example in swagger.yaml for GET /v1/orders/{order_id}.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) GetOrder(ctx context.Context, orderID string) (model.OrderResponse, error) {
	_ = ctx
	_ = orderID

	return model.OrderResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "ACTIVE",
			Mode:   "RESERVATION",
		},
		Payment: &model.PaymentInfo{
			Status: "PAID",
		},
		Charging: &model.ChargingInfo{
			Status: "ACTIVE",
		},
		ConnectorID:   "",
		ConnectorType: "",
		Vehicle: &model.Vehicle{
			Type:  "",
			Model: "",
			Make:  "",
		},
		TrackingURL: "https://track.bluechargenet-aggregator.io/session/SESSION-9876543210",
		ChargingTelemetry: &model.ChargingTelemetry{
			EventTime: "2025-01-27T17:00:00Z",
			Metrics: []model.ChargingMetric{
				{
					Name:     "STATE_OF_CHARGE",
					Value:    62.5,
					UnitCode: "PERCENTAGE",
				},
				{
					Name:     "POWER",
					Value:    18.4,
					UnitCode: "KWH",
				},
				{
					Name:     "ENERGY",
					Value:    10.2,
					UnitCode: "KW",
				},
				{
					Name:     "VOLTAGE",
					Value:    392,
					UnitCode: "VLT",
				},
				{
					Name:     "CURRENT",
					Value:    47,
					UnitCode: "AMP",
				},
				{
					Name:     "SESSION_DURATION",
					Value:    10,
					UnitCode: "min",
				},
			},
		},
	}, nil
}


