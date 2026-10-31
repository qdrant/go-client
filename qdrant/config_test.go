package qdrant_test

import (
	"context"
	"errors"
	"math"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

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

func TestRetryInterceptorWithMaxUintRetries(t *testing.T) {
	calls := 0
	client, err := qdrant.NewGrpcClient(&qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig:            &qdrant.RetryConfig{MaxRetries: math.MaxUint},
		GrpcOptions: []grpc.DialOption{stubCalls(status.Error(codes.InvalidArgument, "stub"), func(context.Context) {
			calls++
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	_, err = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Equal(t, 1, calls)
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
		// The zone separator is percent-encoded in the gRPC target.
		{host: "fe80::1%eth0", want: "[fe80::1%25eth0]:6334"},
	} {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host: tt.host, Port: tt.port, PoolSize: 1, SkipCompatibilityCheck: true,
		})
		require.NoError(t, err, tt.host)
		require.Equal(t, tt.want, client.GetGrpcClient().Conn().Target())
		require.NoError(t, client.Close())
	}
}

// largeReplyServer answers HealthCheck with a title of titleSize bytes.
type largeReplyServer struct {
	qdrant.UnimplementedQdrantServer
	titleSize int
	calls     *atomic.Int32
}

func (s largeReplyServer) HealthCheck(context.Context, *qdrant.HealthCheckRequest) (*qdrant.HealthCheckReply, error) {
	s.calls.Add(1)
	return &qdrant.HealthCheckReply{Title: strings.Repeat("x", s.titleSize)}, nil
}

// newInMemoryClient starts an in-memory server that replies with titleSize bytes
// and returns a client connected to it, plus a count of the calls the server received.
func newInMemoryClient(t *testing.T, titleSize int, cfg *qdrant.Config) (*qdrant.GrpcClient, *atomic.Int32) {
	t.Helper()
	calls := &atomic.Int32{}
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer(grpc.MaxSendMsgSize(math.MaxInt))
	qdrant.RegisterQdrantServer(srv, largeReplyServer{titleSize: titleSize, calls: calls})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	cfg.SkipCompatibilityCheck = true
	cfg.GrpcOptions = append(cfg.GrpcOptions, grpc.WithContextDialer(
		func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) },
	))
	client, err := qdrant.NewGrpcClient(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return client, calls
}

func TestMaxRecvMsgSize(t *testing.T) {
	const mb = 1 << 20

	t.Run("DefaultAcceptsResponseOver4MB", func(t *testing.T) {
		client, _ := newInMemoryClient(t, 5*mb, &qdrant.Config{})
		reply, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Len(t, reply.GetTitle(), 5*mb)
	})

	t.Run("DefaultRejectsResponseOver64MB", func(t *testing.T) {
		client, _ := newInMemoryClient(t, 65*mb, &qdrant.Config{})
		_, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.Equal(t, codes.ResourceExhausted, status.Code(err))
	})

	t.Run("MaxUintIsCappedNotWrapped", func(t *testing.T) {
		client, _ := newInMemoryClient(t, 5*mb, &qdrant.Config{MaxRecvMsgSize: math.MaxUint})
		reply, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Len(t, reply.GetTitle(), 5*mb)
	})

	t.Run("CustomLimitRejectsLargerResponse", func(t *testing.T) {
		client, _ := newInMemoryClient(t, 2*mb, &qdrant.Config{MaxRecvMsgSize: mb})
		_, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.Equal(t, codes.ResourceExhausted, status.Code(err))
	})

	t.Run("GrpcOptionsOverrideConfig", func(t *testing.T) {
		client, _ := newInMemoryClient(t, 2*mb, &qdrant.Config{
			MaxRecvMsgSize: mb,
			GrpcOptions:    []grpc.DialOption{grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(8 * mb))},
		})
		reply, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.NoError(t, err)
		require.Len(t, reply.GetTitle(), 2*mb)
	})
}

func TestRetrySkipsMessageTooLarge(t *testing.T) {
	const mb = 1 << 20

	t.Run("ResponseOverLimitIsNotRetried", func(t *testing.T) {
		client, calls := newInMemoryClient(t, 2*mb, &qdrant.Config{
			MaxRecvMsgSize: mb,
			RetryConfig:    &qdrant.RetryConfig{MaxRetries: 3, BaseBackoff: time.Millisecond},
		})
		_, err := client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
		require.Equal(t, codes.ResourceExhausted, status.Code(err))
		require.Contains(t, status.Convert(err).Message(), "larger than max")
		require.Equal(t, int32(1), calls.Load())
	})

	// grpc-go's size-limit messages, as returned by the invoker.
	for name, msg := range map[string]string{
		"Received":                   "grpc: received message larger than max (2097152 vs. 1048576)",
		"AfterDecompression":         "grpc: message after decompression larger than max (2097152 vs. 1048576)",
		"ReceivedAfterDecompression": "grpc: received message after decompression larger than max 1048576",
		"Send":                       "trying to send message larger than max (2097152 vs. 1048576)",
	} {
		t.Run("NotRetried/"+name, func(t *testing.T) {
			require.Equal(t, 1, retryCalls(t, status.Error(codes.ResourceExhausted, msg)))
		})
	}

	t.Run("ServerResourceExhaustedIsRetried", func(t *testing.T) {
		require.Equal(t, 4, retryCalls(t, status.Error(codes.ResourceExhausted, "rate limit exceeded")))
	})

	t.Run("UnavailableIsRetried", func(t *testing.T) {
		require.Equal(t, 4, retryCalls(t, status.Error(codes.Unavailable, "unavailable")))
	})
}

// retryCalls makes one call that always fails with err, with up to 3 retries,
// and returns how many times the call was attempted.
func retryCalls(t *testing.T, err error) int {
	t.Helper()
	calls := 0
	client, newErr := qdrant.NewGrpcClient(&qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig:            &qdrant.RetryConfig{MaxRetries: 3, BaseBackoff: time.Millisecond},
		GrpcOptions:            []grpc.DialOption{stubCalls(err, func(context.Context) { calls++ })},
	})
	require.NoError(t, newErr)
	t.Cleanup(func() { _ = client.Close() })

	_, _ = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})
	return calls
}
