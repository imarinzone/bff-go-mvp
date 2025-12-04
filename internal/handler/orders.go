package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/orders"
	"bff-go-mvp/internal/httpx"
)

// OrdersHandler handles order-related endpoints.
type OrdersHandler struct {
	service orders.Service
	logger  *zap.Logger
}

func NewOrdersHandler(service orders.Service, logger *zap.Logger) *OrdersHandler {
	return &OrdersHandler{
		service: service,
		logger:  logger,
	}
}

// GetOrder handles GET /v1/orders/{order_id}.
func (h *OrdersHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Required headers: X-Transaction-Id, X-Bpp-Id
	txnID := r.Header.Get("X-Transaction-Id")
	bppID := r.Header.Get("X-Bpp-Id")
	if txnID == "" || bppID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request parameters or missing required fields.")
		return
	}

	vars := mux.Vars(r)
	orderID := vars["order_id"]
	if orderID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Missing order_id in path.")
		return
	}

	resp, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		h.logger.Error("orders service failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	// Echo headers as per swagger
	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}


