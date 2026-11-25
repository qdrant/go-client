package qdrant_test

import (
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/qdrant/go-client/qdrant"
)

// rateLimitCalls returns a dial option that fails each call with ResourceExhausted
// and a retry-after trailer, calling onCall first.
func rateLimitCalls(retryAfter string, onCall func()) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(
		_ context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		onCall()
		for _, opt := range opts {
			if trailer, ok := opt.(grpc.TrailerCallOption); ok {
				*trailer.TrailerAddr = metadata.Pairs("retry-after", retryAfter)
			}
		}
		return status.Error(codes.ResourceExhausted, "stub")
	})
}

func TestRetryWaitsForRetryAfterBeyondMaxBackoff(t *testing.T) {
	calls := 0
	client, err := qdrant.NewGrpcClient(&qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig: &qdrant.RetryConfig{
			MaxRetries:  1,
			BaseBackoff: time.Millisecond,
			MaxBackoff:  time.Millisecond,
		},
		GrpcOptions: []grpc.DialOption{rateLimitCalls("1", func() {
			calls++
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	start := time.Now()
	_, err = client.Qdrant().HealthCheck(context.Background(), &qdrant.HealthCheckRequest{})

	var rateLimitErr *qdrant.QdrantResourceExhaustedError
	require.ErrorAs(t, err, &rateLimitErr)
	require.Equal(t, 2, calls)
	require.GreaterOrEqual(t, time.Since(start), time.Second)
}

func TestRetryAfterPastContextDeadlineReturnsError(t *testing.T) {
	const timeout = 50 * time.Millisecond
	calls := 0
	client, err := qdrant.NewGrpcClient(&qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig: &qdrant.RetryConfig{
			MaxRetries:  1,
			BaseBackoff: time.Millisecond,
			MaxBackoff:  time.Millisecond,
		},
		GrpcOptions: []grpc.DialOption{rateLimitCalls("60", func() {
			calls++
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	start := time.Now()
	_, err = client.Qdrant().HealthCheck(ctx, &qdrant.HealthCheckRequest{})

	var rateLimitErr *qdrant.QdrantResourceExhaustedError
	require.ErrorAs(t, err, &rateLimitErr)
	require.Less(t, time.Since(start), timeout)
	require.Equal(t, 1, calls)
}

func TestLargeRetryAfterDoesNotOverflow(t *testing.T) {
	calls := 0
	client, err := qdrant.NewGrpcClient(&qdrant.Config{
		SkipCompatibilityCheck: true,
		RetryConfig: &qdrant.RetryConfig{
			MaxRetries:  1,
			BaseBackoff: time.Millisecond,
			MaxBackoff:  time.Millisecond,
		},
		// Without capping, seconds * time.Second overflows to a negative wait and the call is retried.
		GrpcOptions: []grpc.DialOption{rateLimitCalls("9223372036854775807", func() {
			calls++
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)
	_, err = client.Qdrant().HealthCheck(ctx, &qdrant.HealthCheckRequest{})

	var rateLimitErr *qdrant.QdrantResourceExhaustedError
	require.ErrorAs(t, err, &rateLimitErr)
	require.Equal(t, 1, calls)
}

type rateLimitHealthServer struct {
	qdrant.UnimplementedQdrantServer
	calls atomic.Int32
}

func (s *rateLimitHealthServer) HealthCheck(
	ctx context.Context, _ *qdrant.HealthCheckRequest,
) (*qdrant.HealthCheckReply, error) {
	call := s.calls.Add(1)
	retryAfter := "1"
	if call > 1 {
		retryAfter = "2"
	}
	if err := grpc.SetTrailer(ctx, metadata.Pairs("retry-after", retryAfter)); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.ResourceExhausted, "rate limited")
}

func TestRetryHonorsRealServerRetryAfterTrailer(t *testing.T) {
	const (
		maxRetries          = 1
		expectedCalls       = 2
		expectedRetryAfterS = 2
	)
	lis, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = lis.Close() })

	srv := grpc.NewServer()
	healthSrv := &rateLimitHealthServer{}
	qdrant.RegisterQdrantServer(srv, healthSrv)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	addr, ok := lis.Addr().(*net.TCPAddr)
	require.True(t, ok)

	client, err := qdrant.NewGrpcClient(&qdrant.Config{
		Host:                   addr.IP.String(),
		Port:                   addr.Port,
		SkipCompatibilityCheck: true,
		RetryConfig: &qdrant.RetryConfig{
			MaxRetries:  maxRetries,
			BaseBackoff: time.Millisecond,
			MaxBackoff:  time.Millisecond,
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	var callerTrailer metadata.MD
	start := time.Now()
	_, err = client.Qdrant().HealthCheck(
		context.Background(),
		&qdrant.HealthCheckRequest{},
		grpc.Trailer(&callerTrailer),
	)

	var rateLimitErr *qdrant.QdrantResourceExhaustedError
	require.ErrorAs(t, err, &rateLimitErr)
	require.Equal(t, expectedRetryAfterS, rateLimitErr.RetryAfterS)
	require.Equal(t, []string{"2"}, callerTrailer.Get("retry-after"))
	require.Equal(t, int32(expectedCalls), healthSrv.calls.Load())
	require.GreaterOrEqual(t, time.Since(start), time.Second)
}
