package support

import (
	"context"

	"bff-go-mvp/internal/model"
)

// Service defines behavior for retrieving support details for an order.
type Service interface {
	GetSupport(ctx context.Context, orderID string) (model.SupportResponse, error)
}


