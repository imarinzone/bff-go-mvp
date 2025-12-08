package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	discoverpb "bff-go-mvp/proto/discover/gen/discover"
)

// DiscoverClient represents a gRPC client for discover service
type DiscoverClient struct {
	conn   *grpc.ClientConn
	client discoverpb.DiscoveryServiceClient
}

// NewDiscoverClient creates a new gRPC discover client
func NewDiscoverClient(address string) (*DiscoverClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := discoverpb.NewDiscoveryServiceClient(conn)

	return &DiscoverClient{
		conn:   conn,
		client: client,
	}, nil
}

// Discover calls the discover service via gRPC using goroutine and channel pattern
func (c *DiscoverClient) Discover(ctx context.Context, req *discoverpb.DiscoverRequest) (*discoverpb.OnDiscoverResponse, error) {
	type result struct {
		resp *discoverpb.OnDiscoverResponse
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		resp, err := c.client.Discover(ctx, req)
		resultChan <- result{resp: resp, err: err}
	}()

	select {
	case res := <-resultChan:
		return res.resp, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close closes the gRPC client connection
func (c *DiscoverClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

