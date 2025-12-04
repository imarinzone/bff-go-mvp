package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bff-go-mvp/internal/domain/payment"
	"bff-go-mvp/internal/httpx"
	"bff-go-mvp/internal/model"
)

// PaymentHandler handles POST /v1/orders/{order_id}/payment requests.
type PaymentHandler struct {
	service payment.Service
	logger  *zap.Logger
}

func NewPaymentHandler(service payment.Service, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

// InitiatePayment handles the payment initiation API.
// @Summary Initiate payment for an order
// @Description Initiates a payment flow for the specified order.
// @Tags Payment
// @Accept json
// @Produce json
// @Param X-Transaction-Id header string true "Unique transaction identifier"
// @Param X-Bpp-Id header string true "Backend provider identifier"
// @Param order_id path string true "Order ID"
// @Param request body object false "Payment initiation payload"
// @Success 200 {object} model.PaymentResponse
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /v1/orders/{order_id}/payment [post]
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Required headers: X-Transaction-Id, X-Bpp-Id
	txnID := r.Header.Get("X-Transaction-Id")
	bppID := r.Header.Get("X-Bpp-Id")
	if txnID == "" || bppID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request parameters or malformed request body.")
		return
	}

	vars := mux.Vars(r)
	orderID := vars["order_id"]
	if orderID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Missing order_id in path.")
		return
	}

	var body map[string]interface{}
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && err.Error() != "EOF" {
			h.logger.Warn("failed to decode payment request body", zap.Error(err))
			httpx.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body.")
			return
		}
	}

	resp, err := h.service.InitiatePayment(r.Context(), orderID, body)
	if err != nil {
		h.logger.Error("payment service failed", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Server error occurred while processing the request.")
		return
	}

	// Echo back headers as per swagger
	w.Header().Set("X-Transaction-Id", txnID)
	w.Header().Set("X-Bpp-Id", bppID)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// keep model types referenced for Swagger annotations
var _ model.PaymentResponse
var _ model.Error
