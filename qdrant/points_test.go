package qdrant_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qdrant/go-client/qdrant"
)

// TestScrollAllRejectsOrderBy ensures ScrollAll does not silently truncate
// ordered scrolls, which have no pagination cursor.
func TestScrollAllRejectsOrderBy(t *testing.T) {
	// The client is never dialed: the iterator must fail before any RPC.
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   "127.0.0.1",
		Port:                   1,
		SkipCompatibilityCheck: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	iterator := client.ScrollAll(context.Background(), &qdrant.ScrollPoints{
		CollectionName: "test",
		OrderBy:        &qdrant.OrderBy{Key: "price"},
	})

	for range 2 {
		points, err := iterator.Next()
		require.ErrorIs(t, err, qdrant.ErrScrollAllOrderByUnsupported)
		require.Nil(t, points)
	}
}
