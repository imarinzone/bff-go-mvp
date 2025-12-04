package payment

import (
	"context"

	"bff-go-mvp/internal/model"
)

// MockService implements Service and returns static data that matches
// the example in swagger.yaml for POST /v1/orders/{order_id}/payment.
type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (s *MockService) InitiatePayment(ctx context.Context, orderID string, body map[string]interface{}) (model.PaymentResponse, error) {
	_ = ctx
	_ = body

	resp := model.PaymentResponse{
		Order: model.OrderInfo{
			ID:     "order-bpp-789012",
			Status: "ACTIVE",
			Mode:   "RESERVATION",
		},
		Amount: model.Amount{
			Value:    128.64,
			Currency: "INR",
		},
		BeneficiaryID: "",
		AcceptedPaymentMethod: []string{
			"BankTransfer",
			"UPI",
			"Wallet",
		},
		PaymentURL: "",
		Validity: &model.Validity{
			StartDate: "2025-01-27T00:00:00Z",
			EndDate:   "2025-04-27T23:59:59Z",
		},
	}

	_ = orderID // in a real impl this would be used

	return resp, nil
}


