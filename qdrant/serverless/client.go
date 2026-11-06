package serverless

import (
	"context"
	"sync/atomic"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

const grpcPort = 443

// Client is a client for Qdrant Serverless.
// Point operations take the same requests as [qdrant.Client].
type Client struct {
	points *qdrant.Client
	// One per pooled connection.
	collections []CollectionsServiceClient
	next        atomic.Uint32
}

// NewClient creates a client from a [qdrant.Config]. It always connects over TLS on port 443,
// ignoring Port and UseTLS, and skips the server version check.
func NewClient(config *qdrant.Config) (*Client, error) {
	cfgCopy := *config
	cfgCopy.Port = grpcPort
	cfgCopy.UseTLS = true
	cfgCopy.SkipCompatibilityCheck = true
	points, err := qdrant.NewClient(&cfgCopy)
	if err != nil {
		return nil, err
	}

	// GetConnection cycles through the pool in order.
	first := points.GetConnection()
	collections := []CollectionsServiceClient{NewCollectionsServiceClient(first)}
	for conn := points.GetConnection(); conn != first; conn = points.GetConnection() {
		collections = append(collections, NewCollectionsServiceClient(conn))
	}
	return &Client{points: points, collections: collections}, nil
}

// Get the low-level client for the serverless collections gRPC service.
// In case of a pool, it returns one of the clients in a round-robin fashion.
func (c *Client) GetCollectionsClient() CollectionsServiceClient {
	idx := c.next.Add(1) - 1
	return c.collections[idx%uint32(len(c.collections))]
}

// Get the low-level client for the points gRPC service.
func (c *Client) GetPointsClient() qdrant.PointsClient {
	return c.points.GetPointsClient()
}

// GetConnection returns one of the underlying gRPC connections from the pool.
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.points.GetConnection()
}

// Close tears down all underlying connections.
func (c *Client) Close() error {
	return c.points.Close()
}

// Creates a collection. GetResult() of the response describes the outcome, e.g. "created".
// An existing collection fails with code AlreadyExists, or is reported through GetResult().
func (c *Client) CreateCollection(ctx context.Context, collectionName string,
	config *CollectionConfig,
) (*CreateCollectionResponse, error) {
	resp, err := c.GetCollectionsClient().CreateCollection(ctx, &CreateCollectionRequest{
		CollectionName: collectionName,
		Config:         config,
	})
	if err != nil {
		return nil, qdrant.NewQdrantError(err, "CreateCollection", collectionName)
	}
	return resp, nil
}

// Deletes a collection and all of its data.
func (c *Client) DeleteCollection(ctx context.Context,
	collectionName string,
) (*DeleteCollectionResponse, error) {
	resp, err := c.GetCollectionsClient().DeleteCollection(ctx, &DeleteCollectionRequest{
		CollectionName: collectionName,
	})
	if err != nil {
		return nil, qdrant.NewQdrantError(err, "DeleteCollection", collectionName)
	}
	return resp, nil
}

// Returns a collection's configuration and point count.
// A missing collection is reported by GetExists(), not as an error.
func (c *Client) GetCollection(ctx context.Context,
	collectionName string,
) (*GetCollectionResponse, error) {
	resp, err := c.GetCollectionsClient().GetCollection(ctx, &GetCollectionRequest{
		CollectionName: collectionName,
	})
	if err != nil {
		return nil, qdrant.NewQdrantError(err, "GetCollection", collectionName)
	}
	return resp, nil
}

// Checks whether a collection exists.
func (c *Client) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	resp, err := c.GetCollection(ctx, collectionName)
	if err != nil {
		return false, err
	}
	return resp.GetExists(), nil
}

// Lists a page of collections. A nil request lists the first page.
// Pass GetNextOffsetToken() as OffsetToken to fetch the next page.
func (c *Client) ListCollections(ctx context.Context,
	request *ListCollectionsRequest,
) (*ListCollectionsResponse, error) {
	if request == nil {
		request = &ListCollectionsRequest{}
	}
	resp, err := c.GetCollectionsClient().ListCollections(ctx, request)
	if err != nil {
		return nil, qdrant.NewQdrantError(err, "ListCollections")
	}
	return resp, nil
}
