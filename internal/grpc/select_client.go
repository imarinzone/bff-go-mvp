package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	selectpb "bff-go-mvp/proto/select/gen/select"
)

// SelectClient represents a gRPC client for select service
type SelectClient struct {
	conn   *grpc.ClientConn
	client selectpb.SelectServiceClient
}

// NewSelectClient creates a new gRPC select client
func NewSelectClient(address string) (*SelectClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := selectpb.NewSelectServiceClient(conn)

	return &SelectClient{
		conn:   conn,
		client: client,
	}, nil
}

// Select calls the select service via gRPC using goroutine and channel pattern
func (c *SelectClient) Select(ctx context.Context, req *selectpb.SelectRequest) (*selectpb.OnSelectResponse, error) {
	type result struct {
		resp *selectpb.OnSelectResponse
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		resp, err := c.client.Select(ctx, req)
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
func (c *SelectClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

