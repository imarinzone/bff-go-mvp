package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"bff-go-mvp/pkg/models"
)

// DiscoveryWorkflowInput contains input parameters for the workflow
type DiscoveryWorkflowInput struct {
	Request        models.DiscoveryRequest
	ServiceAddress string
}

// DiscoveryWorkflow orchestrates the discovery process
func DiscoveryWorkflow(ctx workflow.Context, input DiscoveryWorkflowInput) (models.DiscoveryResponse, error) {
	// Set activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute the activity
	var result models.DiscoveryResponse
	err := workflow.ExecuteActivity(ctx, CallDiscoveryActivity, input.Request, input.ServiceAddress).Get(ctx, &result)
	if err != nil {
		return models.DiscoveryResponse{}, err
	}

	return result, nil
}
