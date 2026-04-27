package qdrant

import (
	"context"
	"math"
	"math/rand/v2"
	"time"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	defaultBaseBackoff = 100 * time.Millisecond
	defaultMaxBackoff  = 5 * time.Second
)

// RetryConfig controls automatic retry behavior for transient gRPC failures.
// When set on Config, a unary interceptor is registered that retries calls
// receiving ResourceExhausted or Unavailable status codes.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	// Zero means no retries.
	MaxRetries uint
	// BaseBackoff is the initial wait duration before the first retry.
	// Defaults to 100ms if zero.
	BaseBackoff time.Duration
	// MaxBackoff caps the backoff duration between retries.
	// Defaults to 5s if zero.
	MaxBackoff time.Duration
}

func (rc *RetryConfig) baseBackoff() time.Duration {
	if rc.BaseBackoff > 0 {
		return rc.BaseBackoff
	}
	return defaultBaseBackoff
}

func (rc *RetryConfig) maxBackoff() time.Duration {
	if rc.MaxBackoff > 0 {
		return rc.MaxBackoff
	}
	return defaultMaxBackoff
}

func (rc *RetryConfig) retryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		for attempt := range rc.MaxRetries + 1 {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			if attempt == rc.MaxRetries {
				break
			}
			if !isRetryable(err) {
				return err
			}
			backoff := rc.backoffDuration(attempt)
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
		return err
	}
}

func isRetryable(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	code := st.Code()
	return code == codes.ResourceExhausted || code == codes.Unavailable
}

const backoffBase = 2.0

// backoffDuration computes an exponential backoff with full jitter.
func (rc *RetryConfig) backoffDuration(attempt uint) time.Duration {
	exp := math.Pow(backoffBase, float64(attempt))
	backoff := min(time.Duration(float64(rc.baseBackoff())*exp), rc.maxBackoff())
	// Full jitter: uniform random in [0, backoff).
	if backoff > 0 {
		backoff = time.Duration(rand.Int64N(int64(backoff)))
	}
	return backoff
}
