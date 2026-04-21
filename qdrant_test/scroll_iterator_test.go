package qdrant_test

import (
	"context"
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestScrollAll(t *testing.T) {
	collectionName := t.Name()
	vectorSize := uint64(4)
	distance := qdrant.Distance_Cosine
	apiKey := "<SCROLL_ITERATOR_TEST>"
	totalPoints := 50
	pageLimit := uint32(10)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	container, err := standaloneQdrant(ctx, apiKey)
	require.NoError(t, err)

	err = container.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6334/tcp")
	require.NoError(t, err)

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   port.Int(),
		APIKey: apiKey,
	})
	require.NoError(t, err)

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: distance,
		}),
	})
	require.NoError(t, err)

	// Insert 50 points in a single upsert.
	points := make([]*qdrant.PointStruct, totalPoints)
	for i := range points {
		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i)),
			Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
			Payload: qdrant.NewValueMap(map[string]any{
				"index": i,
			}),
		}
	}
	wait := true
	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
		Wait:           &wait,
	})
	require.NoError(t, err)

	count, err := client.Count(ctx, &qdrant.CountPoints{
		CollectionName: collectionName,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(totalPoints), count)

	t.Run("IterateAllPages", func(t *testing.T) {
		iter := client.ScrollAll(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
			Limit:          &pageLimit,
		})

		var collected []*qdrant.RetrievedPoint
		pages := 0
		for {
			batch, err := iter.Next()
			if err != nil {
				require.ErrorIs(t, err, io.EOF)
				break
			}
			require.NotEmpty(t, batch)
			collected = append(collected, batch...)
			pages++
		}

		require.Len(t, collected, totalPoints)
		expectedPages := int(math.Ceil(float64(totalPoints) / float64(pageLimit)))
		require.Equal(t, expectedPages, pages)
	})

	t.Run("DoubleEOF", func(t *testing.T) {
		iter := client.ScrollAll(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
			Limit:          &pageLimit,
		})
		// Drain the iterator.
		for {
			_, err := iter.Next()
			if err != nil {
				require.ErrorIs(t, err, io.EOF)
				break
			}
		}
		// Calling Next again should keep returning io.EOF.
		_, err := iter.Next()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("EmptyCollection", func(t *testing.T) {
		emptyCollection := fmt.Sprintf("%s_empty", collectionName)
		err := client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: emptyCollection,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     vectorSize,
				Distance: distance,
			}),
		})
		require.NoError(t, err)

		iter := client.ScrollAll(ctx, &qdrant.ScrollPoints{
			CollectionName: emptyCollection,
			Limit:          &pageLimit,
		})
		batch, err := iter.Next()
		require.ErrorIs(t, err, io.EOF)
		require.Nil(t, batch)

		err = client.DeleteCollection(ctx, emptyCollection)
		require.NoError(t, err)
	})

	err = client.DeleteCollection(ctx, collectionName)
	require.NoError(t, err)
}
