package qdrant

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

// Client is a high-level client for Qdrant.
// It can manage a single connection or a pool of connections, chosen by setting
// PoolSize in the Config.
type Client struct {
	clients       []*GrpcClient
	next          uint32
	healthMonitor *HealthMonitor
}

// NewClient creates a new Qdrant client.
// It checks Config.PoolSize to determine whether to create a single client
// or a pool of clients. If PoolSize > 1, requests are distributed across
// the connections in a round-robin fashion.
func NewClient(config *Config) (*Client, error) {
	// Ensure config is not modified for the caller by cloning.
	cfgCopy := *config
	// Set default values, depending of Cloud bool
	if cfgCopy.PoolSize == 0 {
		if cfgCopy.Cloud {
			cfgCopy.PoolSize = 3
		} else {
			cfgCopy.PoolSize = 1
		}
	}
	// Create the client, with an inner connection pool of go grpc clients
	client := &Client{
		clients: make([]*GrpcClient, 0, cfgCopy.PoolSize),
	}
	// Iterate over the pool size to create the individual client.
	for i := range cfgCopy.PoolSize {
		if i > 0 {
			// In case of a pool, we only want to check compatibility once.
			cfgCopy.SkipCompatibilityCheck = true
		}
		grpcClient, err := NewGrpcClient(&cfgCopy)
		if err != nil {
			// Close already opened clients before returning an error
			client.Close()
			return nil, fmt.Errorf("failed to create client %d in pool: %w", i, err)
		}
		client.clients = append(client.clients, grpcClient)
	}

	// Initialize health monitoring if configured
	if cfgCopy.HealthCheck != nil {
		client.healthMonitor = NewHealthMonitor(client.clients, cfgCopy.HealthCheck)
		client.healthMonitor.Start()
	}

	// Return the client
	return client, nil
}

// Instantiates a new client with the default configuration.
// Connects to localhost:6334 with TLS disabled.
func DefaultClient() (*Client, error) {
	return NewClient(&Config{})
}

// get returns the next GrpcClient from the pool in a round-robin fashion.
// If health monitoring is enabled, it prefers healthy connections.
func (c *Client) get() *GrpcClient {
	if len(c.clients) == 1 {
		return c.clients[0]
	}

	// If health monitoring is enabled, try to get a healthy connection
	if c.healthMonitor != nil {
		healthyIndices := c.healthMonitor.GetHealthyConnections()
		if len(healthyIndices) > 0 {
			// Use round-robin among healthy connections
			idx := atomic.AddUint32(&c.next, 1) - 1
			healthyIdx := healthyIndices[idx%uint32(len(healthyIndices))]
			return c.clients[healthyIdx]
		}
		// If no healthy connections, fall back to regular round-robin
		// This allows for potential recovery attempts
	}

	// Regular round-robin selection
	idx := atomic.AddUint32(&c.next, 1) - 1
	return c.clients[idx%uint32(len(c.clients))]
}

// Get the underlying gRPC client. In case of a pool, it returns one of the clients
// in a round-robin fashion.
func (c *Client) GetGrpcClient() *GrpcClient {
	return c.get()
}

// Get the low-level client for the collections gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/collections_service.proto
func (c *Client) GetCollectionsClient() CollectionsClient {
	return c.get().Collections()
}

// Get the low-level client for the points gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/points_service.proto
func (c *Client) GetPointsClient() PointsClient {
	return c.get().Points()
}

// Get the low-level client for the snapshots gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/snapshots_service.proto
func (c *Client) GetSnapshotsClient() SnapshotsClient {
	return c.get().Snapshots()
}

// Get the low-level client for the Qdrant gRPC service.
// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/qdrant.proto
func (c *Client) GetQdrantClient() QdrantClient {
	return c.get().Qdrant()
}

// GetConnection returns one of the underlying gRPC connections from the pool.
// Useful for debugging or advanced use cases.
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.get().Conn()
}

// Close tears down all underlying connections.
func (c *Client) Close() error {
	// Stop health monitoring first
	if c.healthMonitor != nil {
		c.healthMonitor.Stop()
		c.healthMonitor = nil
	}

	var lastErr error
	for _, client := range c.clients {
		if err := client.Close(); err != nil {
			lastErr = err
		}
	}
	c.clients = nil // Clear the slice
	return lastErr
}

// GetHealthMonitor returns the health monitor instance if health checking is enabled.
// Returns nil if health checking is disabled.
func (c *Client) GetHealthMonitor() *HealthMonitor {
	return c.healthMonitor
}

// GetPoolHealth returns the overall health status of the connection pool.
// Returns nil if health monitoring is disabled.
func (c *Client) GetPoolHealth() *PoolHealth {
	if c.healthMonitor == nil {
		return nil
	}
	health := c.healthMonitor.GetOverallHealth()
	return &health
}

// IsHealthy returns true if at least one connection in the pool is healthy.
// If health monitoring is disabled, always returns true.
func (c *Client) IsHealthy() bool {
	if c.healthMonitor == nil {
		return true // Assume healthy if monitoring is disabled
	}
	return c.healthMonitor.GetOverallHealth().IsHealthy()
}

// WaitForHealthy blocks until at least one connection becomes healthy or the context is cancelled.
// Returns immediately if health monitoring is disabled.
func (c *Client) WaitForHealthy(ctx context.Context) error {
	if c.healthMonitor == nil {
		return nil // No monitoring, assume healthy
	}

	// Check if already healthy
	if c.IsHealthy() {
		return nil
	}

	// Poll for health status
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if c.IsHealthy() {
				return nil
			}
		}
	}
}

// Creates a pointer to a value of any type.
func PtrOf[T any](t T) *T {
	return &t
}
