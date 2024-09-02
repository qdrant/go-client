package qdrant_test

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestPointsClient(t *testing.T) {
	collectionName := t.Name()
	vectorSize := uint64(4)
	distance := qdrant.Distance_Cosine
	apiKey := "<POINTS_TEST>"

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

	testPointID := qdrant.NewID("ed7ac159-d8a7-41fb-9da3-66a14916330f")
	wait := true

	t.Run("UpsertPoints", func(t *testing.T) {
		points := []*qdrant.PointStruct{
			{
				Id:      testPointID,
				Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
			},
		}
		res, err := client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points:         points,
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Test with invalid vector size
		points[0].Vectors = qdrant.NewVectors(0.1, 0.2)
		res, err = client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points:         points,
			Wait:           &wait,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GetPoints", func(t *testing.T) {
		points, err := client.Get(ctx, &qdrant.GetPoints{
			CollectionName: collectionName,
			Ids: []*qdrant.PointId{
				testPointID,
			},
		})
		require.NoError(t, err)
		require.Len(t, points, 1)

		// Test with non-existent point ID
		points, err = client.Get(ctx, &qdrant.GetPoints{
			CollectionName: collectionName,
			Ids: []*qdrant.PointId{
				qdrant.NewIDNum(423),
			},
		})
		require.NoError(t, err)
		require.Empty(t, points)
	})

	t.Run("CountPoints", func(t *testing.T) {
		count, err := client.Count(ctx, &qdrant.CountPoints{
			CollectionName: collectionName,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(1), count)
	})

	t.Run("ScrollPoints", func(t *testing.T) {
		points, err := client.Scroll(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
		})
		require.NoError(t, err)
		require.Len(t, points, 1)
	})

	t.Run("UpdateVectors", func(t *testing.T) {
		points := []*qdrant.PointVectors{
			{
				Id:      testPointID,
				Vectors: qdrant.NewVectors(0.4, 0.5, 0.6, 0.7),
			},
		}
		res, err := client.UpdateVectors(ctx, &qdrant.UpdatePointVectors{
			CollectionName: collectionName,
			Points:         points,
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)

		// Test with invalid vector size
		points[0].Vectors = qdrant.NewVectors(0.1, 0.2)
		res, err = client.UpdateVectors(ctx, &qdrant.UpdatePointVectors{
			CollectionName: collectionName,
			Points:         points,
			Wait:           &wait,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("QueryPoints", func(t *testing.T) {
		res, err := client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(0.1, 0.2, 0.3, 0.4),
		})
		require.NoError(t, err)
		require.Len(t, res, 1)

		// Test with invalid query vector size
		res, err = client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(0.1, 0.2),
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("QueryBatchPoints", func(t *testing.T) {
		res, err := client.QueryBatch(ctx, &qdrant.QueryBatchPoints{
			CollectionName: collectionName,
			QueryPoints: []*qdrant.QueryPoints{
				{
					CollectionName: collectionName,
					Query:          qdrant.NewQuery(0.1, 0.2, 0.3, 0.4),
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("QueryGroups", func(t *testing.T) {
		groups, err := client.QueryGroups(ctx, &qdrant.QueryPointGroups{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(0.1, 0.2, 0.3, 0.4),
			GroupBy:        "key",
		})
		require.NoError(t, err)
		require.Empty(t, groups)
	})

	t.Run("DeleteVectors", func(t *testing.T) {
		res, err := client.DeleteVectors(ctx, &qdrant.DeletePointVectors{
			CollectionName: collectionName,
			Vectors: &qdrant.VectorsSelector{
				// Delete the default/unnamed vector that we're using
				Names: []string{""},
			},
			PointsSelector: qdrant.NewPointsSelector(testPointID),
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("SetPayload", func(t *testing.T) {
		res, err := client.SetPayload(ctx, &qdrant.SetPayloadPoints{
			CollectionName: collectionName,
			PointsSelector: qdrant.NewPointsSelector(testPointID),
			Payload: qdrant.NewValueMap(map[string]any{
				"key":   "value",
				"key_2": 32,
				"key_3": false,
			}),
			Wait: &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("OverwritePayload", func(t *testing.T) {
		res, err := client.OverwritePayload(ctx, &qdrant.SetPayloadPoints{
			CollectionName: collectionName,
			PointsSelector: qdrant.NewPointsSelector(testPointID),
			Payload: qdrant.NewValueMap(map[string]any{
				"key": 10,
			}),
			Wait: &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("DeletePayload", func(t *testing.T) {
		res, err := client.DeletePayload(ctx, &qdrant.DeletePayloadPoints{
			CollectionName: collectionName,
			PointsSelector: qdrant.NewPointsSelector(testPointID),
			Keys:           []string{"key"},
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("ClearPayload", func(t *testing.T) {
		res, err := client.ClearPayload(ctx, &qdrant.ClearPayloadPoints{
			CollectionName: collectionName,
			Points:         qdrant.NewPointsSelector(testPointID),
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("CreateFieldIndex", func(t *testing.T) {
		res, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "key",
			FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("DeleteFieldIndex", func(t *testing.T) {
		res, err := client.DeleteFieldIndex(ctx, &qdrant.DeleteFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "key",
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("UpdateBatch", func(t *testing.T) {
		ops := []*qdrant.PointsUpdateOperation{
			qdrant.NewPointsUpdateDeletePayload(&qdrant.PointsUpdateOperation_DeletePayload{
				Keys:           []string{"key"},
				PointsSelector: qdrant.NewPointsSelector(testPointID),
			}),
		}
		res, err := client.UpdateBatch(ctx, &qdrant.UpdateBatchPoints{
			CollectionName: collectionName,
			Operations:     ops,
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("DeletePoints", func(t *testing.T) {
		res, err := client.Delete(ctx, &qdrant.DeletePoints{
			CollectionName: collectionName,
			Points:         qdrant.NewPointsSelector(testPointID),
			Wait:           &wait,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	err = client.DeleteCollection(ctx, collectionName)
	require.NoError(t, err)
}
