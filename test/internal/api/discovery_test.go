package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"bff-go-mvp/internal/api"
	"bff-go-mvp/pkg/models"
)

func TestDiscoveryHandler_HandleDiscovery(t *testing.T) {
	// Create mock request
	reqBody := models.DiscoveryRequest{
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

	jsonBody, _ := json.Marshal(reqBody)

	// Since we can't easily mock the grpc.NewClient call in the handler,
	// we'll test with the actual implementation which uses a mock response
	// In a real scenario, you might want to refactor to inject the gRPC client
	logger := zap.NewNop()
	handler := api.NewDiscoveryHandler("localhost:50051", logger)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/discovery", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.HandleDiscovery(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp models.DiscoveryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Context.TransactionID, resp.Context.TransactionID)
}

func TestDiscoveryHandler_HandleDiscovery_InvalidMethod(t *testing.T) {
	logger := zap.NewNop()
	handler := api.NewDiscoveryHandler("localhost:50051", logger)

	req := httptest.NewRequest(http.MethodGet, "/discovery", nil)
	w := httptest.NewRecorder()

	handler.HandleDiscovery(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestDiscoveryHandler_HandleDiscovery_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	handler := api.NewDiscoveryHandler("localhost:50051", logger)

	req := httptest.NewRequest(http.MethodPost, "/discovery", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.HandleDiscovery(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
