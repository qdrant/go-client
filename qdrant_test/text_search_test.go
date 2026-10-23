package qdrant_test

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestTextSearch(t *testing.T) {
	collectionName := t.Name()
	apiKey := "<TEXT_SEARCH_TEST>"

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
		Port:   int(port.Num()),
		APIKey: apiKey,
	})
	require.NoError(t, err)

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Dot,
		}),
	})
	require.NoError(t, err)

	wait := true

	t.Run("QueryText", func(t *testing.T) {
		_, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "title",
			FieldType:      qdrant.FieldType_FieldTypeText.Enum(),
			FieldIndexParams: qdrant.NewPayloadIndexParamsText(&qdrant.TextIndexParams{
				Tokenizer: qdrant.TokenizerType_Word,
				Scoring: &qdrant.TextScoringParams{
					Type: qdrant.TextScoringType_Bm25,
				},
			}),
			Wait: &wait,
		})
		require.NoError(t, err)

		_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Wait:           &wait,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewIDNum(1),
					Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
					Payload: qdrant.NewValueMap(map[string]any{"title": "vector search engine"}),
				},
				{
					Id:      qdrant.NewIDNum(2),
					Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
					Payload: qdrant.NewValueMap(map[string]any{"title": "search search search"}),
				},
				{
					Id:      qdrant.NewIDNum(3),
					Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
					Payload: qdrant.NewValueMap(map[string]any{"title": "relational database"}),
				},
			},
		})
		require.NoError(t, err)

		res, err := client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQueryText(&qdrant.TextQuery{Query: "vector search"}),
			Using:          qdrant.PtrOf("title"),
		})
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.Equal(t, uint64(1), res[0].GetId().GetNum())
		require.Equal(t, uint64(2), res[1].GetId().GetNum())

		tuned, err := client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: collectionName,
			Query: qdrant.NewQueryText(qdrant.NewTextQueryBM25("vector search", &qdrant.Bm25Params{
				K: qdrant.PtrOf(float32(2.0)),
				B: qdrant.PtrOf(float32(0.5)),
			})),
			Using: qdrant.PtrOf("title"),
		})
		require.NoError(t, err)
		require.Len(t, tuned, 2)
		require.NotEqual(t, res[0].GetScore(), tuned[0].GetScore())
	})

	t.Run("MatchSubstring", func(t *testing.T) {
		_, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "sku",
			FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
			Wait:           &wait,
		})
		require.NoError(t, err)

		_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Wait:           &wait,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewIDNum(30),
					Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
					Payload: qdrant.NewValueMap(map[string]any{"sku": "AB-1234-XY"}),
				},
				{
					Id:      qdrant.NewIDNum(31),
					Vectors: qdrant.NewVectors(0.5, 0.6, 0.7, 0.8),
					Payload: qdrant.NewValueMap(map[string]any{"sku": "CD-5678-XY"}),
				},
			},
		})
		require.NoError(t, err)

		points, err := client.Scroll(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
			Filter: &qdrant.Filter{
				Must: []*qdrant.Condition{qdrant.NewMatchSubstring("sku", "1234")},
			},
		})
		require.NoError(t, err)
		require.Len(t, points, 1)
		require.Equal(t, uint64(30), points[0].GetId().GetNum())
	})

	t.Run("CodeTokenizer", func(t *testing.T) {
		_, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "source",
			FieldType:      qdrant.FieldType_FieldTypeText.Enum(),
			FieldIndexParams: qdrant.NewPayloadIndexParamsText(&qdrant.TextIndexParams{
				Tokenizer: qdrant.TokenizerType_Code,
			}),
			Wait: &wait,
		})
		require.NoError(t, err)

		_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Wait:           &wait,
			Points: []*qdrant.PointStruct{
				{
					Id:      qdrant.NewIDNum(10),
					Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
					Payload: qdrant.NewValueMap(map[string]any{"source": "fn getUserName() {}"}),
				},
				{
					Id:      qdrant.NewIDNum(11),
					Vectors: qdrant.NewVectors(0.5, 0.6, 0.7, 0.8),
					Payload: qdrant.NewValueMap(map[string]any{"source": "fn deleteFile() {}"}),
				},
			},
		})
		require.NoError(t, err)

		points, err := client.Scroll(ctx, &qdrant.ScrollPoints{
			CollectionName: collectionName,
			Filter: &qdrant.Filter{
				Must: []*qdrant.Condition{qdrant.NewMatchText("source", "get_user_name")},
			},
		})
		require.NoError(t, err)
		require.Len(t, points, 1)
		require.Equal(t, uint64(10), points[0].GetId().GetNum())
	})

	t.Run("NestedFieldIndex", func(t *testing.T) {
		res, err := client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      "offers[].price",
			FieldType:      qdrant.FieldType_FieldTypeFloat.Enum(),
			FieldIndexParams: qdrant.NewPayloadIndexParamsFloat(&qdrant.FloatIndexParams{
				Nested: qdrant.PtrOf("offers"),
			}),
			Wait: &wait,
		})
		require.NoError(t, err)
		require.Equal(t, qdrant.UpdateStatus_Completed, res.GetStatus())
	})

	err = client.DeleteCollection(ctx, collectionName)
	require.NoError(t, err)
}
