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

func buildRouter() http.Handler {
	logger := zap.NewNop()
	cfg := config.Load()
	return router.New(cfg, logger)
}

func TestOrdersLifecycle_EstimateCancel_Success(t *testing.T) {
	r := buildRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/orders/order-123/cancel?activity=test", nil)
	req.Header.Set("X-Transaction-Id", "txn-1")
	req.Header.Set("X-Bpp-Id", "bpp-1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.CancelEstimateResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "order-bpp-789012", resp.Order.ID)
}

func TestOrdersLifecycle_Cancel_Success(t *testing.T) {
	r := buildRouter()

	body := map[string]interface{}{"reason": "user_cancel"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders/order-123/cancel", bytes.NewReader(bodyBytes))
	req.Header.Set("X-Transaction-Id", "txn-1")
	req.Header.Set("X-Bpp-Id", "bpp-1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp model.CancelResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "CANCELLED", resp.Order.Status)
}

func TestOrdersLifecycle_Start_Stop_Success(t *testing.T) {
	r := buildRouter()

	// Start
	startReq := httptest.NewRequest(http.MethodPut, "/v1/orders/order-123/start", nil)
	startReq.Header.Set("X-Transaction-Id", "txn-1")
	startReq.Header.Set("X-Bpp-Id", "bpp-1")
	startW := httptest.NewRecorder()
	r.ServeHTTP(startW, startReq)
	assert.Equal(t, http.StatusAccepted, startW.Code)

	// Stop
	stopBody := model.StopChargingRequest{ReasonCode: "USER", Message: "Stop now"}
	stopBytes, _ := json.Marshal(stopBody)
	stopReq := httptest.NewRequest(http.MethodPut, "/v1/orders/order-123/stop", bytes.NewReader(stopBytes))
	stopReq.Header.Set("X-Transaction-Id", "txn-1")
	stopReq.Header.Set("X-Bpp-Id", "bpp-1")
	stopW := httptest.NewRecorder()
	r.ServeHTTP(stopW, stopReq)
	assert.Equal(t, http.StatusOK, stopW.Code)

	var stopResp model.StopChargingResponse
	err := json.Unmarshal(stopW.Body.Bytes(), &stopResp)
	assert.NoError(t, err)
	assert.Equal(t, "COMPLETED", stopResp.Order.Status)
}
