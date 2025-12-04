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

func buildTestRouter() http.Handler {
	logger := zap.NewNop()
	cfg := config.Load()
	return router.New(cfg, logger)
}

func TestFeedbackHandler_SetOrderRating_Success(t *testing.T) {
	r := buildTestRouter()

	reqBody := model.RatingRequest{
		Value: 5,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/orders/order-123/rating", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transaction-Id", "txn-1")
	req.Header.Set("X-Bpp-Id", "bpp-1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp model.RatingResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "order-bpp-789012", resp.Order.ID)
	assert.NotNil(t, resp.FeedbackForm)
}

func TestSupportHandler_GetOrderSupport_Success(t *testing.T) {
	r := buildTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/orders/order-123/support", nil)
	req.Header.Set("X-Transaction-Id", "txn-1")
	req.Header.Set("X-Bpp-Id", "bpp-1")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.SupportResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "BlueCharge Support Team", resp.Name)
	assert.Contains(t, resp.Channels, "phone")
}


