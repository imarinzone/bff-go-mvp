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

func TestSearchHandler_Success(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()

	r := router.New(cfg, logger)

	reqBody := model.SearchRequest{
		GeoCoordinates: []float64{12.9716, 77.5946},
		DistanceMeters: 5000,
	}

	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/search?page=1&per_page=20", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("X-Transaction-Id"))

	var resp model.SearchResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, 1, resp.Page)
}

func TestSearchHandler_InvalidOneOf(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()
	r := router.New(cfg, logger)

	// Neither evse_id nor geo_coordinates provided -> bad request
	reqBody := model.SearchRequest{}
	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/search", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchHandler_DiscoverRequestFormat(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.Load()

	r := router.New(cfg, logger)

	// Test with the full discover request format (context + message)
	discoverReq := model.DiscoverRequest{
		Context: model.DiscoverContext{
			TransactionID: "test-e2e-001",
			Version:       "2.0.0",
			Action:        "discover",
			Domain:        "nic2004:52110",
			BapID:         "bap.example.com",
			BapURI:        "https://bap.example.com",
			MessageID:     "msg-001",
			Timestamp:     "2025-12-05T17:50:00+05:30",
			TTL:           "PT30M",
		},
		Message: model.DiscoverMessage{
			Geometry: &model.Geometry{
				Type:        "Point",
				Coordinates: []float64{77.2090, 28.6139},
			},
			DistanceMeters: 5000.0,
		},
	}

	bodyBytes, err := json.Marshal(discoverReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/search?page=1&per_page=20", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("X-Transaction-Id"))

	var resp model.SearchResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, resp.Total, 0)
	assert.Equal(t, 1, resp.Page)
}
