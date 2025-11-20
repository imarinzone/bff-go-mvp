package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/pkg/models"
)

// DiscoveryHandler handles discovery API requests
type DiscoveryHandler struct {
	grpcClient     *grpc.Client
	logger         *zap.Logger
	serviceAddress string
}

// NewDiscoveryHandler creates a new discovery handler
func NewDiscoveryHandler(serviceAddress string, logger *zap.Logger) *DiscoveryHandler {
	return &DiscoveryHandler{
		serviceAddress: serviceAddress,
		logger:         logger,
	}
}

// HandleDiscovery handles POST /discovery requests
func (h *DiscoveryHandler) HandleDiscovery(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method for discovery endpoint",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.DiscoveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	h.logger.Info("Processing discovery request",
		zap.String("transaction_id", req.Context.TransactionID),
		zap.String("message_id", req.Context.MessageID),
	)

	// Create gRPC client and call discovery service
	grpcClient := grpc.NewClient(h.serviceAddress)
	defer grpcClient.Close()

	result, err := grpcClient.CallDiscoveryService(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to call discovery service",
			zap.Error(err),
			zap.String("transaction_id", req.Context.TransactionID),
		)
		http.Error(w, fmt.Sprintf("Failed to call discovery service: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Discovery request completed successfully",
		zap.String("transaction_id", req.Context.TransactionID),
		zap.String("message_id", req.Context.MessageID),
	)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("Failed to encode response",
			zap.Error(err),
			zap.String("transaction_id", req.Context.TransactionID),
		)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
