package qdrant

import (
	"google.golang.org/grpc"
)

// Lower level client for Qdrant gRPC API.
type GrpcClient struct {
	conn *grpc.ClientConn
	// Qdrant service interface
	// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/qdrant.proto
	qdrant QdrantClient
	// Collections service interface
	// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/collections_service.proto
	collections CollectionsClient
	// Points service interface
	// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/points_service.proto
	points PointsClient
	//  Snapshots service interface
	// https://github.com/qdrant/qdrant/blob/master/lib/api/src/grpc/proto/snapshots_service.proto
	snapshots SnapshotsClient
}

// Create a new gRPC client with default configuration.
func NewDefaultGrpcClient() (*GrpcClient, error) {
	return NewGrpcClient(&Config{})
}

// Create a new gRPC client with custom configuration.
func NewGrpcClient(config *Config) (*GrpcClient, error) {
	// We append config.GrpcOptions in the end
	// so that user's explicit options take precedence
	config.GrpcOptions = append([]grpc.DialOption{
		config.getTransportCreds(),
		config.getAPIKeyInterceptor(),
	}, config.GrpcOptions...)

	conn, err := grpc.NewClient(config.getAddr(), config.GrpcOptions...)

	if err != nil {
		return nil, err
	}

	return NewGrpcClientFromConn(conn), nil
}

// Create a new gRPC client from existing connection.
func NewGrpcClientFromConn(conn *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{
		conn:        conn,
		qdrant:      NewQdrantClient(conn),
		points:      NewPointsClient(conn),
		collections: NewCollectionsClient(conn),
		snapshots:   NewSnapshotsClient(conn),
	}
}

// Get the underlying gRPC connection.
func (c *GrpcClient) Conn() *grpc.ClientConn {
	return c.conn
}

// Get the Qdrant service interface.
func (c *GrpcClient) Qdrant() QdrantClient {
	return c.qdrant
}

// Get the Collections service interface.
func (c *GrpcClient) Points() PointsClient {
	return c.points
}

// Get the Points service interface.
func (c *GrpcClient) Collections() CollectionsClient {
	return c.collections
}

// Get the Snapshots service interface.
func (c *GrpcClient) Snapshots() SnapshotsClient {
	return c.snapshots
}

// Tears down the *grpc.ClientConn and all underlying connections.
func (c *GrpcClient) Close() error {
	return c.conn.Close()
}
