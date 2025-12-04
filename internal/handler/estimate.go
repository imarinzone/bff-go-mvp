package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/estimate"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// EstimateHandler handles POST /v1/estimate requests.
type EstimateHandler struct {
	service estimate.Service
	logger  *zap.Logger
}

func NewEstimateHandler(service estimate.Service, logger *zap.Logger) *EstimateHandler {
	return &EstimateHandler{
		service: service,
		logger:  logger,
	}
}

// GetEstimates handles the estimate API.
func (h *EstimateHandler) GetEstimates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Required header in swagger: X-Transaction-Id
	txnID := r.Header.Get("X-Transaction-Id")
	if txnID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Missing required header X-Transaction-Id")
		return
	}

	var req model.EstimateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode estimate request", zap.Error(err))
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	// Basic required field validation per swagger: evse_id, vehicle, connector_id
	if req.EvseID == "" || req.ConnectorID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request parameters or missing required fields.")
		return
	}

	resp, err := h.service.Estimate(r.Context(), req)
	if err != nil {
		h.logger.Error("estimate service failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	// Echo back transaction id and set a mock BPP id
	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", "mock-bpp-id")
	httpx.WriteJSON(w, http.StatusOK, resp)
}


