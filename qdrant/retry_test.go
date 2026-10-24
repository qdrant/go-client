package qdrant_test

import (
	"context"
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

func TestRetryAfterWaitStopsAtContextDeadline(t *testing.T) {
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

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), time.Second)
	require.Equal(t, 1, calls)
}
