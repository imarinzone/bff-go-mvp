package orders

import (
	"context"

	"bff-go-mvp/internal/model"
)

// LifecycleService defines operations for order lifecycle: start, stop, cancel.
type LifecycleService interface {
	EstimateCancel(ctx context.Context, orderID, activity, cancelReason, cancelCode string) (model.CancelEstimateResponse, error)
	Cancel(ctx context.Context, orderID string, body map[string]interface{}) (model.CancelResponse, error)
	EstimateStop(ctx context.Context, orderID, activity string) (model.StopEstimateResponse, error)
	Stop(ctx context.Context, orderID string, req model.StopChargingRequest) (model.StopChargingResponse, error)
	Start(ctx context.Context, orderID string, req model.StartChargingRequest) (model.StartChargingResponse, error)
}

// MockLifecycleService implements LifecycleService using static example data.
type MockLifecycleService struct{}

func NewMockLifecycleService() *MockLifecycleService {
	return &MockLifecycleService{}
}

func (s *MockLifecycleService) EstimateCancel(ctx context.Context, orderID, activity, cancelReason, cancelCode string) (model.CancelEstimateResponse, error) {
	_ = ctx
	_ = orderID
	_ = activity
	_ = cancelReason
	_ = cancelCode

	return model.CancelEstimateResponse{
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
		Validity: &model.Validity{
			StartDate: "2025-01-27T00:00:00Z",
			EndDate:   "2025-04-27T23:59:59Z",
		},
		PriceComponents: []model.PriceComponentString{
			{
				Type:        "PAID",
				Value:       "400.00",
				Currency:    "INR",
				Description: "Base price",
			},
			{
				Type:        "FEE",
				Value:       "30.00",
				Currency:    "INR",
				Description: "Cancellation charges",
			},
		},
	}, nil
}

func (s *MockLifecycleService) Cancel(ctx context.Context, orderID string, body map[string]interface{}) (model.CancelResponse, error) {
	_ = ctx
	_ = orderID
	_ = body

	return model.CancelResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "CANCELLED",
			Mode:   "RESERVATION",
		},
		Payment: &model.PaymentInfo{
			Status: "PAID",
		},
		Charging: &model.ChargingInfo{
			Status: "ACTIVE",
		},
		PriceComponents: []model.PriceComponentString{
			{
				Type:        "FEE",
				Value:       "30.00",
				Currency:    "INR",
				Description: "Cancellation charges",
			},
			{
				Type:        "REFUND",
				Value:       "-300.00",
				Currency:    "INR",
				Description: "Cancellation refund",
			},
		},
	}, nil
}

func (s *MockLifecycleService) EstimateStop(ctx context.Context, orderID, activity string) (model.StopEstimateResponse, error) {
	_ = ctx
	_ = orderID
	_ = activity

	return model.StopEstimateResponse{
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
		Validity: &model.Validity{
			StartDate: "2025-01-27T00:00:00Z",
			EndDate:   "2025-04-27T23:59:59Z",
		},
		PriceComponents: []model.PriceComponentString{
			{
				Type:        "PAID",
				Value:       "400.00",
				Currency:    "INR",
				Description: "Base price",
			},
			{
				Type:        "FEE",
				Value:       "30.00",
				Currency:    "INR",
				Description: "Cancellation charges",
			},
			{
				Type:        "REFUND",
				Value:       "-300.00",
				Currency:    "INR",
				Description: "Cancellation refund",
			},
		},
	}, nil
}

func (s *MockLifecycleService) Stop(ctx context.Context, orderID string, req model.StopChargingRequest) (model.StopChargingResponse, error) {
	_ = ctx
	_ = orderID
	_ = req

	return model.StopChargingResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "COMPLETED",
			Mode:   "RESERVATION",
		},
		Payment: &model.PaymentInfo{
			Status: "PAID",
		},
		Charging: &model.ChargingInfo{
			Status: "COMPLETED",
		},
		Validity: &model.Validity{
			StartDate: "2025-01-27T00:00:00Z",
			EndDate:   "2025-04-27T23:59:59Z",
		},
		PriceComponents: []model.PriceComponentFlexible{
			{
				Type:        "BASE",
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
				Type:        "REFUND",
				Value:       "-300.00",
				Currency:    "INR",
				Description: "Cancellation refund",
			},
		},
	}, nil
}

func (s *MockLifecycleService) Start(ctx context.Context, orderID string, req model.StartChargingRequest) (model.StartChargingResponse, error) {
	_ = ctx
	_ = orderID
	_ = req

	return model.StartChargingResponse{
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
	}, nil
}


