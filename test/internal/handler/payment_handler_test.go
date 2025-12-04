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
	"bff-go-mvp/internal/router"
)

func TestPaymentHandler_Success(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	body := map[string]interface{}{
		"dummy": "value",
	}
	bodyBytes, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders/order-123/payment", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transaction-Id", "txn-abc")
	req.Header.Set("X-Bpp-Id", "bpp-xyz")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "txn-abc", w.Header().Get("X-Transaction-Id"))
	assert.Equal(t, "bpp-xyz", w.Header().Get("X-Bpp-Id"))
}

func TestPaymentHandler_MissingHeaders(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders/order-123/payment", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}


