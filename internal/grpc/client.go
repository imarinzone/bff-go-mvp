package grpc

import (
	"context"
	"fmt"

	"bff-go-mvp/pkg/models"
)

// Client represents a gRPC client for discovery service
type Client struct {
	serviceAddress string
}

// NewClient creates a new gRPC client
func NewClient(serviceAddress string) *Client {
	return &Client{
		serviceAddress: serviceAddress,
	}
}

// CallDiscoveryService calls the discovery service via gRPC and returns a mock response
// In a real implementation, this would make an actual gRPC call
func (c *Client) CallDiscoveryService(ctx context.Context, req *models.DiscoveryRequest) (*models.DiscoveryResponse, error) {
	// Mock implementation - returns the request with some modifications
	// In production, this would:
	// 1. Convert models.DiscoveryRequest to protobuf DiscoveryRequest
	// 2. Make gRPC call to the service
	// 3. Convert protobuf DiscoveryResponse to models.DiscoveryResponse

	// For now, return a mock response based on the request
	response := &models.DiscoveryResponse{
		Context: models.Context{
			Version:       req.Context.Version,
			Action:        "on_discover",
			Domain:        req.Context.Domain,
			Location:      req.Context.Location,
			BapID:         req.Context.BapID,
			BapURI:        req.Context.BapURI,
			BppID:         req.Context.BppID,
			BppURI:        req.Context.BppURI,
			TransactionID: req.Context.TransactionID,
			MessageID:     fmt.Sprintf("resp-%s", req.Context.MessageID),
			Timestamp:     req.Context.Timestamp,
			TTL:           req.Context.TTL,
			SchemaContext: req.Context.SchemaContext,
		},
		Message: models.Message{
			Catalogs: req.Message.Catalogs, // Echo back the catalogs for now
		},
	}

	return response, nil
}

// Close closes the gRPC client connection
func (c *Client) Close() error {
	// In a real implementation, this would close the gRPC connection
	return nil
}

