package temporal

import (
	"context"

	"bff-go-mvp/internal/grpc"
	"bff-go-mvp/pkg/models"
)

// CallDiscoveryActivity calls the gRPC discovery service
func CallDiscoveryActivity(ctx context.Context, req models.DiscoveryRequest, serviceAddress string) (models.DiscoveryResponse, error) {
	// Create gRPC client
	grpcClient := grpc.NewClient(serviceAddress)
	defer grpcClient.Close()

	// Call the discovery service
	response, err := grpcClient.CallDiscoveryService(ctx, &req)
	if err != nil {
		return models.DiscoveryResponse{}, err
	}

	return *response, nil
}
