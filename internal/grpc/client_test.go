package grpc

import (
	"context"
	"testing"

	"bff-go-mvp/pkg/models"
)

func TestClient_CallDiscoveryService(t *testing.T) {
	client := NewClient("localhost:50051")

	req := &models.DiscoveryRequest{
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

	ctx := context.Background()
	resp, err := client.CallDiscoveryService(ctx, req)
	if err != nil {
		t.Fatalf("CallDiscoveryService failed: %v", err)
	}

	if resp.Context.TransactionID != req.Context.TransactionID {
		t.Errorf("Expected transaction ID %s, got %s", req.Context.TransactionID, resp.Context.TransactionID)
	}

	if resp.Context.MessageID == req.Context.MessageID {
		t.Errorf("Expected response message ID to be different from request")
	}
}

func TestClient_Close(t *testing.T) {
	client := NewClient("localhost:50051")
	err := client.Close()
	if err != nil {
		t.Errorf("Close() should not return an error: %v", err)
	}
}

