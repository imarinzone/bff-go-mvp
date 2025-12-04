package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/feedback"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// FeedbackHandler handles rating/feedback endpoints.
type FeedbackHandler struct {
	service feedback.Service
	logger  *zap.Logger
}

func NewFeedbackHandler(service feedback.Service, logger *zap.Logger) *FeedbackHandler {
	return &FeedbackHandler{
		service: service,
		logger:  logger,
	}
}

// SetOrderRating handles POST /v1/orders/{order_id}/rating.
func (h *FeedbackHandler) SetOrderRating(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	txnID := r.Header.Get("X-Transaction-Id")
	bppID := r.Header.Get("X-Bpp-Id")
	if txnID == "" || bppID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid rating value or malformed request.")
		return
	}

	vars := mux.Vars(r)
	orderID := vars["order_id"]
	if orderID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Missing order_id in path.")
		return
	}

	var req model.RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode rating request", zap.Error(err))
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid rating value or malformed request.")
		return
	}

	if req.Value < 1 || req.Value > 5 {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid rating value or malformed request.")
		return
	}

	resp, err := h.service.SetRating(r.Context(), orderID, req)
	if err != nil {
		h.logger.Error("set rating failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", bppID)
	httpx.WriteJSON(w, http.StatusCreated, resp)
}


