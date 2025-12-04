package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/orders"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// OrdersLifecycleHandler handles start/stop/cancel related endpoints.
type OrdersLifecycleHandler struct {
	service orders.LifecycleService
	logger  *zap.Logger
}

func NewOrdersLifecycleHandler(service orders.LifecycleService, logger *zap.Logger) *OrdersLifecycleHandler {
	return &OrdersLifecycleHandler{
		service: service,
		logger:  logger,
	}
}

// EstimateCancel handles GET /v1/orders/{order_id}/cancel.
func (h *OrdersLifecycleHandler) EstimateCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	txnID, bppID, orderID, ok := h.validateHeadersAndOrderID(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	activity := q.Get("activity")
	cancelReason := q.Get("cancel_reason")
	cancelCode := q.Get("cancel_code")

	resp, err := h.service.EstimateCancel(r.Context(), orderID, activity, cancelReason, cancelCode)
	if err != nil {
		h.logger.Error("estimate cancel failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	h.writeStandardHeaders(w, txnID, bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Cancel handles POST /v1/orders/{order_id}/cancel.
func (h *OrdersLifecycleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	txnID, bppID, orderID, ok := h.validateHeadersAndOrderID(w, r)
	if !ok {
		return
	}

	var body map[string]interface{}
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err.Error() != "EOF" {
			h.logger.Warn("failed to decode cancel request body", zap.Error(err))
			httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body.")
			return
		}
	}

	resp, err := h.service.Cancel(r.Context(), orderID, body)
	if err != nil {
		h.logger.Error("cancel failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	h.writeStandardHeaders(w, txnID, bppID)
	httpx.WriteJSON(w, http.StatusAccepted, resp)
}

// EstimateStop handles GET /v1/orders/{order_id}/stop.
func (h *OrdersLifecycleHandler) EstimateStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	txnID, bppID, orderID, ok := h.validateHeadersAndOrderID(w, r)
	if !ok {
		return
	}

	activity := r.URL.Query().Get("activity")

	resp, err := h.service.EstimateStop(r.Context(), orderID, activity)
	if err != nil {
		h.logger.Error("estimate stop failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	h.writeStandardHeaders(w, txnID, bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// StopCharging handles PUT /v1/orders/{order_id}/stop.
func (h *OrdersLifecycleHandler) StopCharging(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	txnID, bppID, orderID, ok := h.validateHeadersAndOrderID(w, r)
	if !ok {
		return
	}

	var req model.StopChargingRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			h.logger.Warn("failed to decode stop request", zap.Error(err))
			httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body.")
			return
		}
	}

	resp, err := h.service.Stop(r.Context(), orderID, req)
	if err != nil {
		h.logger.Error("stop charging failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	h.writeStandardHeaders(w, txnID, bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// StartCharging handles PUT /v1/orders/{order_id}/start.
func (h *OrdersLifecycleHandler) StartCharging(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	txnID, bppID, orderID, ok := h.validateHeadersAndOrderID(w, r)
	if !ok {
		return
	}

	var req model.StartChargingRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			h.logger.Warn("failed to decode start request", zap.Error(err))
			httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body.")
			return
		}
	}

	resp, err := h.service.Start(r.Context(), orderID, req)
	if err != nil {
		h.logger.Error("start charging failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	h.writeStandardHeaders(w, txnID, bppID)
	httpx.WriteJSON(w, http.StatusAccepted, resp)
}

func (h *OrdersLifecycleHandler) validateHeadersAndOrderID(w http.ResponseWriter, r *http.Request) (txnID, bppID, orderID string, ok bool) {
	txnID = r.Header.Get("X-Transaction-Id")
	bppID = r.Header.Get("X-Bpp-Id")
	if txnID == "" || bppID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request parameters or missing required fields.")
		return "", "", "", false
	}

	vars := mux.Vars(r)
	orderID = vars["order_id"]
	if orderID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Missing order_id in path.")
		return "", "", "", false
	}

	return txnID, bppID, orderID, true
}

func (h *OrdersLifecycleHandler) writeStandardHeaders(w http.ResponseWriter, txnID, bppID string) {
	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", bppID)
}


