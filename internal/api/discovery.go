package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"bff-go-mvp/internal/temporal"
	"bff-go-mvp/pkg/models"
)

// TemporalWorkflowExecutor is an interface for executing Temporal workflows
type TemporalWorkflowExecutor interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

// DiscoveryHandler handles discovery API requests
type DiscoveryHandler struct {
	workflowExecutor TemporalWorkflowExecutor
	serviceAddress   string
	taskQueue        string
	logger           *zap.Logger
}

// NewDiscoveryHandler creates a new discovery handler
func NewDiscoveryHandler(workflowExecutor TemporalWorkflowExecutor, serviceAddress, taskQueue string, logger *zap.Logger) *DiscoveryHandler {
	return &DiscoveryHandler{
		workflowExecutor: workflowExecutor,
		serviceAddress:   serviceAddress,
		taskQueue:        taskQueue,
		logger:           logger,
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

	// Generate workflow ID
	workflowID := fmt.Sprintf("discovery-%s-%s", req.Context.TransactionID, req.Context.MessageID)

	h.logger.Info("Processing discovery request",
		zap.String("workflow_id", workflowID),
		zap.String("transaction_id", req.Context.TransactionID),
		zap.String("message_id", req.Context.MessageID),
	)

	// Prepare workflow input
	workflowInput := temporal.DiscoveryWorkflowInput{
		Request:        req,
		ServiceAddress: h.serviceAddress,
	}

	// Execute Temporal workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: h.taskQueue,
	}

	we, err := h.workflowExecutor.ExecuteWorkflow(r.Context(), workflowOptions, temporal.DiscoveryWorkflow, workflowInput)
	if err != nil {
		h.logger.Error("Failed to start workflow",
			zap.Error(err),
			zap.String("workflow_id", workflowID),
		)
		http.Error(w, fmt.Sprintf("Failed to start workflow: %v", err), http.StatusInternalServerError)
		return
	}

	// Get workflow result
	var result models.DiscoveryResponse
	if err := we.Get(r.Context(), &result); err != nil {
		h.logger.Error("Workflow execution failed",
			zap.Error(err),
			zap.String("workflow_id", workflowID),
		)
		http.Error(w, fmt.Sprintf("Workflow execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Discovery request completed successfully",
		zap.String("workflow_id", workflowID),
		zap.String("transaction_id", req.Context.TransactionID),
	)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("Failed to encode response",
			zap.Error(err),
			zap.String("workflow_id", workflowID),
		)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
