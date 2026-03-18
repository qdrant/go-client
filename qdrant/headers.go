package qdrant

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// WithHeader returns a new context that sends the key and value as a header
// with every request made using that context.
//
//	ctx = qdrant.WithHeader(ctx, "x-request-id", "abc-123")
//	client.Query(ctx, ...)
func WithHeader(ctx context.Context, key, value string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, value)
}

// WithHeaders returns a new context that sends all entries of headers with every
// request made using that context.
//
//	ctx = qdrant.WithHeaders(ctx, map[string]string{"x-request-id": "abc-123"})
//	client.Query(ctx, ...)
func WithHeaders(ctx context.Context, headers map[string]string) context.Context {
	if len(headers) == 0 {
		return ctx
	}
	pairs := make([]string, 0)
	for k, v := range headers {
		pairs = append(pairs, k, v)
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}
