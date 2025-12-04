package orders

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines the behavior for fetching order details.
type Service interface {
	GetOrder(ctx context.Context, orderID string) (model.OrderResponse, error)
}


