package qdrant_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/qdrant/go-client/qdrant"
)

// metadataRecorder is a Qdrant service that records the metadata of the last HealthCheck.
type metadataRecorder struct {
	qdrant.UnimplementedQdrantServer
	last metadata.MD
}

func (m *metadataRecorder) HealthCheck(
	ctx context.Context,
	_ *qdrant.HealthCheckRequest,
) (*qdrant.HealthCheckReply, error) {
	m.last, _ = metadata.FromIncomingContext(ctx)
	return &qdrant.HealthCheckReply{}, nil
}

func startMetadataRecorder(t *testing.T) (*metadataRecorder, string, int) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	recorder := &metadataRecorder{}
	server := grpc.NewServer()
	qdrant.RegisterQdrantServer(server, recorder)
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)

	addr, ok := lis.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return recorder, addr.IP.String(), addr.Port
}

func TestMetadataInterceptor(t *testing.T) {
	t.Run("SendsAPIKeyAndHeaders", func(t *testing.T) {
		recorder, host, port := startMetadataRecorder(t)

		client, err := qdrant.NewGrpcClient(&qdrant.Config{
			Host:                   host,
			Port:                   port,
			APIKey:                 "secret",
			Headers:                map[string]string{"x-request-id": "abc-123"},
			SkipCompatibilityCheck: true,
		})
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Close() })

		_, err = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Equal(t, []string{"secret"}, recorder.last.Get("api-key"))
		require.Equal(t, []string{"abc-123"}, recorder.last.Get("x-request-id"))
	})

	t.Run("NoAPIKeyNoHeaders", func(t *testing.T) {
		recorder, host, port := startMetadataRecorder(t)

		client, err := qdrant.NewGrpcClient(&qdrant.Config{
			Host:                   host,
			Port:                   port,
			SkipCompatibilityCheck: true,
		})
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Close() })

		_, err = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Empty(t, recorder.last.Get("api-key"))
	})

	t.Run("ConfigChangesAfterCreationAreIgnored", func(t *testing.T) {
		recorder, host, port := startMetadataRecorder(t)

		config := &qdrant.Config{
			Host:                   host,
			Port:                   port,
			APIKey:                 "secret",
			Headers:                map[string]string{"x-request-id": "abc-123"},
			SkipCompatibilityCheck: true,
		}
		client, err := qdrant.NewGrpcClient(config)
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Close() })

		config.APIKey = "changed"
		config.Headers["x-request-id"] = "changed"

		_, err = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Equal(t, []string{"secret"}, recorder.last.Get("api-key"))
		require.Equal(t, []string{"abc-123"}, recorder.last.Get("x-request-id"))
	})
}
