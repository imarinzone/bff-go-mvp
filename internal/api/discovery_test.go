package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"bff-go-mvp/pkg/models"
)

// MockWorkflowExecutor is a mock implementation of TemporalWorkflowExecutor
type MockWorkflowExecutor struct {
	mock.Mock
}

func (m *MockWorkflowExecutor) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	args_m := m.Called(ctx, options, workflow, args)
	return args_m.Get(0).(client.WorkflowRun), args_m.Error(1)
}

// MockWorkflowRun is a mock implementation of WorkflowRun
type MockWorkflowRun struct {
	mock.Mock
	client.WorkflowRun
}

func (m *MockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	args := m.Called(ctx, valuePtr)
	// Set the value if provided
	if resp, ok := valuePtr.(*models.DiscoveryResponse); ok {
		*resp = args.Get(0).(models.DiscoveryResponse)
	}
	return args.Error(1)
}

func (m *MockWorkflowRun) GetID() string {
	return "test-workflow-id"
}

func (m *MockWorkflowRun) GetRunID() string {
	return "test-run-id"
}

func TestDiscoveryHandler_HandleDiscovery(t *testing.T) {
	// Create mock request
	reqBody := models.DiscoveryRequest{
		Context: models.Context{
			Version:       "1.0.0",
			Action:        "on_discover",
			Domain:        "mobility",
			TransactionID: "test-txn-123",
			MessageID:     "test-msg-456",
		},
		Message: models.Message{
			Catalogs: []models.Catalog{},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Create mock workflow executor
	mockExecutor := new(MockWorkflowExecutor)
	mockWorkflowRun := new(MockWorkflowRun)

	expectedResponse := models.DiscoveryResponse{
		Context: reqBody.Context,
		Message: reqBody.Message,
	}

	mockWorkflowRun.On("Get", mock.Anything, mock.Anything).Return(expectedResponse, nil)
	mockExecutor.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)

	// Create handler with no-op logger for testing
	logger := zap.NewNop()
	handler := NewDiscoveryHandler(mockExecutor, "localhost:50051", "TEST_QUEUE", logger)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/discovery", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.HandleDiscovery(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp models.DiscoveryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Context.TransactionID, resp.Context.TransactionID)

	mockExecutor.AssertExpectations(t)
	mockWorkflowRun.AssertExpectations(t)
}

func TestDiscoveryHandler_HandleDiscovery_InvalidMethod(t *testing.T) {
	mockExecutor := new(MockWorkflowExecutor)
	logger := zap.NewNop()
	handler := NewDiscoveryHandler(mockExecutor, "localhost:50051", "TEST_QUEUE", logger)

	req := httptest.NewRequest(http.MethodGet, "/discovery", nil)
	w := httptest.NewRecorder()

	handler.HandleDiscovery(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestDiscoveryHandler_HandleDiscovery_InvalidJSON(t *testing.T) {
	mockExecutor := new(MockWorkflowExecutor)
	logger := zap.NewNop()
	handler := NewDiscoveryHandler(mockExecutor, "localhost:50051", "TEST_QUEUE", logger)

	req := httptest.NewRequest(http.MethodPost, "/discovery", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.HandleDiscovery(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
