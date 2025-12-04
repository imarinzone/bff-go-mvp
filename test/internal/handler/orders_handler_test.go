package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/model"
	"bff-go-mvp/internal/router"
	"encoding/json"
)

func TestOrdersHandler_GetOrder_Success(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/v1/orders/order-123", nil)
	req.Header.Set("X-Transaction-Id", "txn-1")
	req.Header.Set("X-Bpp-Id", "bpp-1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "txn-1", w.Header().Get("X-Transaction-Id"))
	assert.Equal(t, "bpp-1", w.Header().Get("X-Bpp-Id"))

	var resp model.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "order-bpp-789012", resp.Order.ID)
	assert.Equal(t, "https://track.bluechargenet-aggregator.io/session/SESSION-9876543210", resp.TrackingURL)
	assert.NotNil(t, resp.ChargingTelemetry)
	assert.GreaterOrEqual(t, len(resp.ChargingTelemetry.Metrics), 1)
}

func TestOrdersHandler_GetOrder_MissingHeaders(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/v1/orders/order-123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}


