package qdrant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/qdrant/go-client/qdrant"
)

func stubPointsCalls(err error, onCall func(ctx context.Context, opts []grpc.CallOption)) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(
		ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		onCall(ctx, opts)
		return err
	})
}

func TestPointOperations_PropagateContextHeaders(t *testing.T) {
	var sent metadata.MD
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(ctx context.Context, _ []grpc.CallOption) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()
	ctx = qdrant.WithHeader(ctx, "x-request-id", "req-123")
	ctx = qdrant.WithHeaders(ctx, map[string]string{"x-tenant-id": "tenant-a"})

	_, _ = client.Upsert(ctx, &qdrant.UpsertPoints{CollectionName: "test"})
	require.Equal(t, []string{"req-123"}, sent.Get("x-request-id"))
	require.Equal(t, []string{"tenant-a"}, sent.Get("x-tenant-id"))

	_, _ = client.Get(ctx, &qdrant.GetPoints{CollectionName: "test"})
	require.Equal(t, []string{"req-123"}, sent.Get("x-request-id"))
	require.Equal(t, []string{"tenant-a"}, sent.Get("x-tenant-id"))

	_, _ = client.Query(ctx, &qdrant.QueryPoints{CollectionName: "test"})
	require.Equal(t, []string{"req-123"}, sent.Get("x-request-id"))
	require.Equal(t, []string{"tenant-a"}, sent.Get("x-tenant-id"))

	_, _ = client.Delete(ctx, &qdrant.DeletePoints{CollectionName: "test"})
	require.Equal(t, []string{"req-123"}, sent.Get("x-request-id"))
	require.Equal(t, []string{"tenant-a"}, sent.Get("x-tenant-id"))
}

func TestPointOperations_WithOptions(t *testing.T) {
	var sent metadata.MD
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(ctx context.Context, _ []grpc.CallOption) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()

	// WithPointHeader
	_, _ = client.Upsert(ctx, &qdrant.UpsertPoints{CollectionName: "test"},
		qdrant.WithPointHeader("x-point-header", "point-val"))
	require.Equal(t, []string{"point-val"}, sent.Get("x-point-header"))

	// WithPointHeaders
	_, _ = client.Query(ctx, &qdrant.QueryPoints{CollectionName: "test"},
		qdrant.WithPointHeaders(map[string]string{
			"x-multi-one": "val1",
			"x-multi-two": "val2",
		}))
	require.Equal(t, []string{"val1"}, sent.Get("x-multi-one"))
	require.Equal(t, []string{"val2"}, sent.Get("x-multi-two"))
}

func TestPointOperations_CombinedHeaders(t *testing.T) {
	var sent metadata.MD
	cfg := &qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		Headers:                map[string]string{"x-config-header": "config-val"},
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(ctx context.Context, _ []grpc.CallOption) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	}
	baseClient, err := qdrant.NewClient(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = baseClient.Close() })

	scopedClient := baseClient.WithHeader("x-client-header", "client-val")

	ctx := qdrant.WithHeader(context.Background(), "x-ctx-header", "ctx-val")

	_, _ = scopedClient.Upsert(ctx, &qdrant.UpsertPoints{CollectionName: "test"},
		qdrant.WithPointHeader("x-opt-header", "opt-val"))

	require.Equal(t, []string{"config-val"}, sent.Get("x-config-header"))
	require.Equal(t, []string{"client-val"}, sent.Get("x-client-header"))
	require.Equal(t, []string{"ctx-val"}, sent.Get("x-ctx-header"))
	require.Equal(t, []string{"opt-val"}, sent.Get("x-opt-header"))
}

func TestClient_WithHeaderAndWithHeadersImmutability(t *testing.T) {
	var sent metadata.MD
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(ctx context.Context, _ []grpc.CallOption) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	scopedClient := client.WithHeader("x-header-1", "val1").WithHeaders(map[string]string{"x-header-2": "val2"})

	// Call using scoped client
	_, _ = scopedClient.Get(context.Background(), &qdrant.GetPoints{CollectionName: "test"})
	require.Equal(t, []string{"val1"}, sent.Get("x-header-1"))
	require.Equal(t, []string{"val2"}, sent.Get("x-header-2"))

	// Call using original client - should not have scoped headers
	sent = nil
	_, _ = client.Get(context.Background(), &qdrant.GetPoints{CollectionName: "test"})
	require.Empty(t, sent.Get("x-header-1"))
	require.Empty(t, sent.Get("x-header-2"))
}

func TestClient_OutgoingContext(t *testing.T) {
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions:            []grpc.DialOption{stubCalls(errors.New("stub"), func(context.Context) {})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	scopedClient := client.WithHeader("x-client-k", "client-v")
	ctx := scopedClient.OutgoingContext(context.Background())

	md, ok := metadata.FromOutgoingContext(ctx)
	require.True(t, ok)
	require.Equal(t, []string{"client-v"}, md.Get("x-client-k"))
}

func TestPointOperations_WithPointGrpcOption(t *testing.T) {
	var capturedOpts []grpc.CallOption
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(_ context.Context, opts []grpc.CallOption) {
			capturedOpts = opts
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	var respHeader metadata.MD
	callOpt := grpc.Header(&respHeader)

	_, _ = client.Upsert(context.Background(), &qdrant.UpsertPoints{CollectionName: "test"},
		qdrant.WithPointGrpcOption(callOpt))

	require.Contains(t, capturedOpts, callOpt)
}

func TestScrollAll_PropagatesPointOptions(t *testing.T) {
	var sent metadata.MD
	client, err := qdrant.NewClient(&qdrant.Config{
		PoolSize:               1,
		SkipCompatibilityCheck: true,
		GrpcOptions: []grpc.DialOption{stubPointsCalls(errors.New("stub"), func(ctx context.Context, _ []grpc.CallOption) {
			sent, _ = metadata.FromOutgoingContext(ctx)
		})},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	it := client.ScrollAll(context.Background(), &qdrant.ScrollPoints{CollectionName: "test"},
		qdrant.WithPointHeader("x-scroll-header", "scroll-val"))

	_, _ = it.Next()

	require.Equal(t, []string{"scroll-val"}, sent.Get("x-scroll-header"))
}
