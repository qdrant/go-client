package qdrant

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func TestRetryInterceptor_MaxUintDoesNotSkip(t *testing.T) {
	rc := &RetryConfig{
		MaxRetries:  math.MaxUint,
		BaseBackoff: 1 * time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
	}

	interceptor := rc.retryInterceptor()
	called := 0
	expectedErr := errors.New("non-retryable error")

	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		called++
		return expectedErr
	}

	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
	assert.Equal(t, 1, called, "invoker should be called even when MaxRetries is math.MaxUint")
	assert.Equal(t, expectedErr, err)
}

func TestRetryInterceptor_Success(t *testing.T) {
	rc := &RetryConfig{
		MaxRetries: 3,
	}

	interceptor := rc.retryInterceptor()
	called := 0

	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		called++
		return nil
	}

	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
	require.NoError(t, err)
	assert.Equal(t, 1, called)
}

func TestRetryInterceptor_RetriesExhausted(t *testing.T) {
	rc := &RetryConfig{
		MaxRetries:  2,
		BaseBackoff: 1 * time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
	}

	interceptor := rc.retryInterceptor()
	called := 0
	retryableErr := status.Error(codes.Unavailable, "service unavailable")

	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		called++
		return retryableErr
	}

	err := interceptor(context.Background(), "/test.Service/Method", nil, nil, nil, invoker)
	assert.Equal(t, 3, called, "should attempt initial call + 2 retries")
	assert.Equal(t, retryableErr, err)
}

func TestRetryInterceptor_ContextCanceled(t *testing.T) {
	rc := &RetryConfig{
		MaxRetries:  5,
		BaseBackoff: 50 * time.Millisecond,
		MaxBackoff:  100 * time.Millisecond,
	}

	interceptor := rc.retryInterceptor()
	ctx, cancel := context.WithCancel(context.Background())
	called := 0
	retryableErr := status.Error(codes.Unavailable, "service unavailable")

	invoker := func(c context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		called++
		cancel()
		return retryableErr
	}

	err := interceptor(ctx, "/test.Service/Method", nil, nil, nil, invoker)
	assert.Equal(t, 1, called)
	assert.Equal(t, context.Canceled, err)
}
