package estimate

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data that matches
// the example in swagger.yaml for POST /v1/estimate.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) Estimate(ctx context.Context, req model.EstimateRequest) (model.EstimateResponse, error) {
	_ = ctx

	resp := model.EstimateResponse{
		Order: model.OrderInfo{
			ID:     "1231208-id",
			Mode:   "reservation",
			Status: "quoted_price",
		},
		Amount: model.Amount{
			Value:    128.64,
			Currency: "INR",
		},
		DurationInMinutes:          "15",
		PercentageOfBatteryCharged: "80",
		Energy: &model.Energy{
			Value: 30,
			Unit:  "kWh",
		},
		Validity: &model.Validity{
			StartDate: "2025-01-27T00:00:00Z",
			EndDate:   "2025-04-27T23:59:59Z",
		},
		PriceComponents: []model.PriceComponent{
			{
				Type:        "UNIT",
				Value:       100,
				Currency:    "INR",
				Description: "Base charging session cost (100 INR)",
			},
			{
				Type:        "SURCHARGE",
				Value:       20,
				Currency:    "INR",
				Description: "Surge price (20%)",
			},
			{
				Type:        "DISCOUNT",
				Value:       -15,
				Currency:    "INR",
				Description: "Offer discount (15%)",
			},
			{
				Type:        "FEE",
				Value:       10,
				Currency:    "INR",
				Description: "Service fee",
			},
			{
				Type:        "FEE",
				Value:       13.64,
				Currency:    "INR",
				Description: "Overcharge estimation",
			},
			{
				Type:        "Pending Payment",
				Value:       0.64,
				Currency:    "INR",
				Description: "Pending payment",
			},
		},
		Cancellation: &model.CancellationPolicy{
			Fee: &model.CancellationFee{
				Percentage: "30",
			},
			ExternalRef: &model.ExternalRef{
				MIMEType: "text/html",
				URL:      "https://example-company.com/charge/tnc.html",
			},
		},
	}

	return resp, nil
}


