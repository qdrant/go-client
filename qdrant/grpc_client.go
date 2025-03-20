package qdrant

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

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
	clientVersion := getClientVersion()
	config.GrpcOptions = append([]grpc.DialOption{
		config.getTransportCreds(),
		config.getAPIKeyInterceptor(),
		config.getRateLimitInterceptor(),
		grpc.WithUserAgent(fmt.Sprintf("go-client/%s", clientVersion)),
	}, config.GrpcOptions...)

	conn, err := grpc.NewClient(config.getAddr(), config.GrpcOptions...)

	if err != nil {
		return nil, err
	}

	newGrpcClientFromConn := NewGrpcClientFromConn(conn)

	if !config.SkipCompatibilityCheck {
		serverVersion := getServerVersion(newGrpcClientFromConn)
		logger := slog.Default()
		if serverVersion == unknownVersion {
			logger.Warn("Failed to obtain server version. " +
				"Unable to check client-server compatibility. " +
				"Set SkipCompatibilityCheck=true to skip version check.")
		} else if !IsCompatible(clientVersion, serverVersion) {
			logger.Warn("Client version is not compatible with server version. "+
				"Major versions should match and minor version difference must not exceed 1. "+
				"Set SkipCompatibilityCheck=true to skip version check.",
				"clientVersion", clientVersion, "serverVersion", serverVersion)
		}
	}

	return newGrpcClientFromConn, nil
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

func getClientVersion() string {
	packageName := "github.com/qdrant/go-client"
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Version}}", packageName)
	output, err := cmd.Output()
	if err != nil {
		return unknownVersion
	}
	return strings.TrimSpace(string(output))
}
