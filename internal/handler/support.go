package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/support"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// SupportHandler handles GET /v1/orders/{order_id}/support.
type SupportHandler struct {
	service support.Service
	logger  *zap.Logger
}

func NewSupportHandler(service support.Service, logger *zap.Logger) *SupportHandler {
	return &SupportHandler{
		service: service,
		logger:  logger,
	}
}

// GetOrderSupport handles GET /v1/orders/{order_id}/support.
// @Summary Get support contact information for an order
// @Description Returns support contact channels and metadata for a specific order.
// @Tags Support
// @Accept json
// @Produce json
// @Param X-Transaction-Id header string true "Unique transaction identifier"
// @Param X-Bpp-Id header string true "Backend provider identifier"
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.SupportResponse
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /v1/orders/{order_id}/support [get]
func (h *SupportHandler) GetOrderSupport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

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

	resp, err := h.service.GetSupport(r.Context(), orderID)
	if err != nil {
		h.logger.Error("get support failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// keep model types referenced for Swagger annotations
var _ model.SupportResponse
var _ model.Error
