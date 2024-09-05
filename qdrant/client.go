package qdrant

import (
	"google.golang.org/grpc"
)

// High-level client for interacting with a Qdrant server.
type Client struct {
	grpcClient *GrpcClient
}

// Instantiates a new client with the given configuration.
func NewClient(config *Config) (*Client, error) {
	grpcClient, err := NewGrpcClient(config)
	if err != nil {
		return nil, err
	}
	return NewClientFromGrpc(grpcClient), nil
}

// Instantiates a new client with the default configuration.
// Connects to localhost:6334 with TLS disabled.
func DefaultClient() (*Client, error) {
	grpcClient, err := NewDefaultGrpcClient()
	if err != nil {
		return nil, err
	}
	return NewClientFromGrpc(grpcClient), err
}

// Instantiates a new client from an existing gRPC client.
func NewClientFromGrpc(grpcClient *GrpcClient) *Client {
	return &Client{
		grpcClient,
	}
}

// Get the underlying gRPC client.
func (c *Client) GetGrpcClient() *GrpcClient {
	return c.grpcClient
}

// Get the low-level client for the collections gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/collections_service.proto
func (c *Client) GetCollectionsClient() CollectionsClient {
	return c.GetGrpcClient().Collections()
}

// Get the low-level client for the points gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/points_service.proto
func (c *Client) GetPointsClient() PointsClient {
	return c.GetGrpcClient().Points()
}

// Get the low-level client for the snapshots gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/snapshots_service.proto
func (c *Client) GetSnapshotsClient() SnapshotsClient {
	return c.GetGrpcClient().Snapshots()
}

// Get the low-level client for the Qdrant gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/qdrant.proto
func (c *Client) GetQdrantClient() QdrantClient {
	return c.GetGrpcClient().Qdrant()
}

// Get the underlying *grpc.ClientConn.
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.GetGrpcClient().Conn()
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}
