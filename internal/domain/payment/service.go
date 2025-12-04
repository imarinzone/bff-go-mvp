package payment

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines the behavior for the payment initiation use case.
type Service interface {
	InitiatePayment(ctx context.Context, orderID string, body map[string]interface{}) (model.PaymentResponse, error)
}


