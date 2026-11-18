package qdrant

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// Client is a high-level client for Qdrant.
// It can manage a single connection or a pool of connections, chosen by setting
// PoolSize in the Config.
type Client struct {
	clients   []*GrpcClient
	next      uint32
	closeOnce sync.Once
	headers   map[string]string
}

// NewClient creates a new Qdrant client.
// It checks Config.PoolSize to determine whether to create a single client
// or a pool of clients. If PoolSize > 1, requests are distributed across
// the connections in a round-robin fashion.
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errNilConfig
	}
	// Ensure config is not modified for the caller by cloning.
	cfgCopy := *config
	cfgCopy.Headers = maps.Clone(config.Headers)
	if cfgCopy.PoolSize == 0 {
		cfgCopy.PoolSize = 3
	}
	// Create the client, with an inner connection pool of go grpc clients
	client := &Client{
		clients: make([]*GrpcClient, 0, cfgCopy.PoolSize),
		headers: make(map[string]string),
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
	// Return the client
	return client, nil
}

// Instantiates a new client with the default configuration.
// Connects to localhost:6334 with TLS disabled.
func DefaultClient() (*Client, error) {
	return NewClient(&Config{})
}

// get returns the next GrpcClient from the pool in a round-robin fashion.
func (c *Client) get() *GrpcClient {
	if len(c.clients) == 1 {
		return c.clients[0]
	}
	// Atomically increment and wrap around the counter
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
// If several connections fail to close, all errors are returned joined.
func (c *Client) Close() error {
	var errs []error
	c.closeOnce.Do(func() {
		for _, client := range c.clients {
			if err := client.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	})
	return errors.Join(errs...)
}

// WithHeader returns a shallow copy of the Client configured with an additional
// outgoing gRPC metadata header sent on all requests made through this client.
func (c *Client) WithHeader(key, value string) *Client {
	if c == nil {
		return nil
	}
	newHeaders := make(map[string]string, len(c.headers)+1)
	maps.Copy(newHeaders, c.headers)
	newHeaders[key] = value
	return &Client{
		clients: c.clients,
		next:    c.next,
		headers: newHeaders,
	}
}

// WithHeaders returns a shallow copy of the Client configured with additional
// outgoing gRPC metadata headers sent on all requests made through this client.
func (c *Client) WithHeaders(headers map[string]string) *Client {
	if c == nil {
		return nil
	}
	newHeaders := make(map[string]string, len(c.headers)+len(headers))
	maps.Copy(newHeaders, c.headers)
	maps.Copy(newHeaders, headers)
	return &Client{
		clients: c.clients,
		next:    c.next,
		headers: newHeaders,
	}
}

// OutgoingContext returns a context enriched with the client's custom metadata headers,
// as well as any existing outgoing metadata.
func (c *Client) OutgoingContext(ctx context.Context) context.Context {
	return c.withOutgoingContext(ctx)
}

func (c *Client) withOutgoingContext(ctx context.Context) context.Context {
	if c == nil || len(c.headers) == 0 {
		return ctx
	}
	return WithHeaders(ctx, c.headers)
}

// Creates a pointer to a value of any type.
func PtrOf[T any](t T) *T {
	return &t
}
