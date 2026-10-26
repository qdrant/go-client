package qdrant

import "google.golang.org/grpc"

// RetryInterceptor exports rc.retryInterceptor for tests in qdrant_test.
func (rc *RetryConfig) RetryInterceptor() grpc.UnaryClientInterceptor {
	return rc.retryInterceptor()
}
