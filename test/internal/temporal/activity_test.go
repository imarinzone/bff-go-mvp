package temporal_test

import (
	"context"
	"testing"

	"bff-go-mvp/internal/temporal"
	"bff-go-mvp/pkg/models"
)

func TestCallDiscoveryActivity(t *testing.T) {
	req := models.DiscoveryRequest{
		Context: models.Context{
			Version:       "1.0.0",
			Action:        "on_discover",
			Domain:        "mobility",
			TransactionID: "test-txn-123",
			MessageID:     "test-msg-456",
		},
		Message: models.Message{
			Catalogs: []models.Catalog{},
		},
	}

	serviceAddress := "localhost:50051"
	ctx := context.Background()

	resp, err := temporal.CallDiscoveryActivity(ctx, req, serviceAddress)
	if err != nil {
		t.Fatalf("CallDiscoveryActivity failed: %v", err)
	}

	if resp.Context.TransactionID != req.Context.TransactionID {
		t.Errorf("Expected transaction ID %s, got %s", req.Context.TransactionID, resp.Context.TransactionID)
	}
}
