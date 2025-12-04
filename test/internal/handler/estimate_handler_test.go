package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/model"
	"bff-go-mvp/internal/router"
)

func TestEstimateHandler_Success(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	reqBody := model.EstimateRequest{
		EvseID:      "evse-123",
		ConnectorID: "connector-456",
		Vehicle: model.Vehicle{
			Make:  "TestMake",
			Model: "TestModel",
			Type:  "4-wheeler",
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/estimate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transaction-Id", "txn-123")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "txn-123", w.Header().Get("X-Transaction-Id"))
	assert.NotEmpty(t, w.Header().Get("X-Bpp-Id"))

	var resp model.EstimateResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "1231208-id", resp.Order.ID)
}

func TestEstimateHandler_MissingTransactionID(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	reqBody := model.EstimateRequest{
		EvseID:      "evse-123",
		ConnectorID: "connector-456",
		Vehicle: model.Vehicle{
			Type: "4-wheeler",
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/estimate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}


