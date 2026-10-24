package qdrant_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/qdrant/go-client/qdrant"
)

// stubCalls returns a dial option that runs after the client's own interceptors.
// It passes each call's context to onCall and fails the call with err, so no server is needed.
func stubCalls(err error, onCall func(ctx context.Context)) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(
		ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ grpc.UnaryInvoker, _ ...grpc.CallOption,
	) error {
		onCall(ctx)
		return err
	})
}

func TestNewGrpcClientIgnoresLaterAPIKeyAndHeaderChanges(t *testing.T) {
	var sent metadata.MD
	cfg := &qdrant.Config{
		APIKey:                 "secret",
		Headers:                map[string]string{"x-id": "1"},
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubCalls(errors.New("stub"), func(ctx context.Context) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	}
	client, err := qdrant.NewGrpcClient(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	cfg.APIKey = "changed"
	cfg.Headers["x-id"] = "changed"
	_, _ = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})

	require.Equal(t, []string{"secret"}, sent.Get("api-key"))
	require.Equal(t, []string{"1"}, sent.Get("x-id"))
}

func TestNewGrpcClientIgnoresLaterRetryConfigChanges(t *testing.T) {
	const maxRetries = 2
	calls := 0
	cfg := &qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig:            &qdrant.RetryConfig{MaxRetries: maxRetries, BaseBackoff: time.Millisecond},
		GrpcOptions: []grpc.DialOption{stubCalls(status.Error(codes.Unavailable, "stub"), func(context.Context) {
			calls++
		})},
	}
	client, err := qdrant.NewGrpcClient(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	cfg.RetryConfig.MaxRetries = 0
	_, _ = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})

	require.Equal(t, maxRetries+1, calls)
}

func TestNilConfigReturnsError(t *testing.T) {
	client, err := qdrant.NewClient(nil)
	require.Error(t, err)
	require.Nil(t, client)

	grpcClient, err := qdrant.NewGrpcClient(nil)
	require.Error(t, err)
	require.Nil(t, grpcClient)
}

func TestClientHostAddress(t *testing.T) {
	for _, tt := range []struct {
		host string
		port int
		want string
	}{
		{host: "", want: "localhost:6334"},
		{host: "qdrant.example.com", port: 6335, want: "qdrant.example.com:6335"},
		{host: "127.0.0.1", want: "127.0.0.1:6334"},
		{host: "::1", want: "[::1]:6334"},
		{host: "[::1]", want: "[::1]:6334"},
	} {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host: tt.host, Port: tt.port, PoolSize: 1, SkipCompatibilityCheck: true,
		})
		require.NoError(t, err, tt.host)
		require.Equal(t, tt.want, client.GetGrpcClient().Conn().Target())
		require.NoError(t, client.Close())
	}
}
