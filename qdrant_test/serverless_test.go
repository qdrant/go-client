package qdrant_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/qdrant/go-client/qdrant/serverless"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Connects to the serverless space at QDRANT_SERVERLESS_URL, e.g.
// https://serverless.example.cloud.qdrant.io. Skips the test if it is not configured.
func serverlessClient(t *testing.T) *serverless.Client {
	t.Helper()
	rawURL := os.Getenv("QDRANT_SERVERLESS_URL")
	apiKey := os.Getenv("QDRANT_SERVERLESS_API_KEY")
	if rawURL == "" || apiKey == "" {
		t.Skip("QDRANT_SERVERLESS_URL and QDRANT_SERVERLESS_API_KEY are not set")
	}
	parsed, err := url.Parse(rawURL)
	require.NoError(t, err)

	client, err := serverless.NewClient(&qdrant.Config{
		Host:   parsed.Hostname(),
		APIKey: apiKey,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})
	return client
}

func serverlessCollection(t *testing.T, client *serverless.Client, config *serverless.CollectionConfig) string {
	t.Helper()
	name := fmt.Sprintf("go-client-test-%d", time.Now().UnixNano())
	_, err := client.CreateCollection(t.Context(), name, config)
	require.NoError(t, err)
	t.Cleanup(func() {
		// Bounded retries: a flaky endpoint shouldn't hang the run or leak the collection.
		var err error
		for range 3 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_, err = client.DeleteCollection(ctx, name)
			cancel()
			if err == nil {
				return
			}
			time.Sleep(time.Second)
		}
		t.Errorf("could not delete collection %s: %v", name, err)
	})
	return name
}

func TestServerlessPort(t *testing.T) {
	client, err := serverless.NewClient(&qdrant.Config{
		Host: "serverless.example.cloud.qdrant.io",
		Port: 6334,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})
	require.Equal(t, "serverless.example.cloud.qdrant.io:443", client.GetConnection().Target())
}

func TestServerlessCollections(t *testing.T) {
	client := serverlessClient(t)
	ctx := t.Context()

	config := &serverless.CollectionConfig{
		DenseVectors: map[string]*serverless.DenseVectorConfig{
			"": {
				Size:          4,
				Distance:      serverless.Distance_COSINE,
				PrecisionTier: serverless.PrecisionTier_FULL.Enum(),
			},
		},
		SparseVectors: map[string]*serverless.SparseVectorConfig{
			"bm25": {UseIdf: true},
		},
		PayloadIndexes: map[string]*serverless.PayloadIndexConfig{
			"color":   serverless.NewPayloadIndexConfigKeyword(nil),
			"age":     serverless.NewPayloadIndexConfigInteger(&serverless.IntegerIndex{Range: qdrant.PtrOf(false)}),
			"price":   serverless.NewPayloadIndexConfigFloat(),
			"user_id": serverless.NewPayloadIndexConfigUUID(),
			"created": serverless.NewPayloadIndexConfigDatetime(),
			"geo":     serverless.NewPayloadIndexConfigGeo(),
			"active":  serverless.NewPayloadIndexConfigBool(),
			"title": serverless.NewPayloadIndexConfigText(&serverless.TextIndex{
				Tokenizer: serverless.Tokenizer_WORD.Enum(),
				Stopwords: &serverless.StopwordsSet{Languages: []string{"english"}},
				Stemmer:   serverless.NewStemmingAlgorithmSnowball("english"),
			}),
			"body": serverless.NewPayloadIndexConfigText(&serverless.TextIndex{
				Stemmer: serverless.NewStemmingAlgorithmDisabled(),
			}),
		},
	}
	name := serverlessCollection(t, client, config)

	t.Run("CreateCollectionAlreadyExists", func(t *testing.T) {
		// The API may report an existing collection as an error or through the result.
		resp, err := client.CreateCollection(ctx, name, config)
		if err != nil {
			require.Equal(t, codes.AlreadyExists, status.Code(err))
		} else {
			require.NotEqual(t, "created", resp.GetResult())
		}
	})

	t.Run("GetCollection", func(t *testing.T) {
		resp, err := client.GetCollection(ctx, name)
		require.NoError(t, err)
		require.True(t, resp.GetExists())
		require.True(t, proto.Equal(config, resp.GetConfig()), "got: %v", resp.GetConfig())

		exists, err := client.CollectionExists(ctx, name)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("GetMissingCollection", func(t *testing.T) {
		resp, err := client.GetCollection(ctx, name+"-missing")
		require.NoError(t, err)
		require.False(t, resp.GetExists())
	})

	t.Run("ListCollections", func(t *testing.T) {
		var names []string
		request := &serverless.ListCollectionsRequest{Limit: qdrant.PtrOf(uint32(100))}
		for page := 0; ; page++ {
			require.Less(t, page, 100, "pagination did not end")
			resp, err := client.ListCollections(ctx, request)
			require.NoError(t, err)
			for _, collection := range resp.GetCollections() {
				names = append(names, collection.GetCollectionName())
			}
			token := resp.GetNextOffsetToken()
			if token == "" {
				break
			}
			request.OffsetToken = &token
		}
		require.Contains(t, names, name)
	})
}

func TestServerlessPoints(t *testing.T) {
	client := serverlessClient(t)
	ctx := t.Context()
	// Reads can lag a few seconds behind writes.
	eventually := func(condition func() bool) {
		require.Eventually(t, condition, time.Minute, 250*time.Millisecond)
	}
	name := serverlessCollection(t, client, &serverless.CollectionConfig{
		DenseVectors: map[string]*serverless.DenseVectorConfig{
			"": {Size: 4, Distance: serverless.Distance_DOT},
		},
		PayloadIndexes: map[string]*serverless.PayloadIndexConfig{
			"color": serverless.NewPayloadIndexConfigKeyword(nil),
		},
	})

	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: name,
		Wait:           qdrant.PtrOf(true),
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(1, 0, 0, 0),
				Payload: qdrant.NewValueMap(map[string]any{"color": "red"}),
			},
			{
				Id:      qdrant.NewIDNum(2),
				Vectors: qdrant.NewVectors(0, 1, 0, 0),
				Payload: qdrant.NewValueMap(map[string]any{"color": "blue"}),
			},
		},
	})
	require.NoError(t, err)

	eventually(func() bool {
		count, err := client.Count(ctx, &qdrant.CountPoints{
			CollectionName: name,
			Exact:          qdrant.PtrOf(true),
			Filter:         &qdrant.Filter{Must: []*qdrant.Condition{qdrant.NewMatch("color", "red")}},
		})
		return err == nil && count == 1
	})

	eventually(func() bool {
		hits, err := client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: name,
			Query:          qdrant.NewQuery(1, 0, 0, 0),
		})
		return err == nil && len(hits) == 2 && hits[0].GetId().GetNum() == 1
	})

	result, err := client.SetPayload(ctx, &qdrant.SetPayloadPoints{
		CollectionName: name,
		Wait:           qdrant.PtrOf(true),
		PointsSelector: qdrant.NewPointsSelector(qdrant.NewIDNum(1)),
		Payload:        qdrant.NewValueMap(map[string]any{"size": 10}),
	})
	require.NoError(t, err)
	require.Equal(t, qdrant.UpdateStatus_Completed, result.GetStatus())

	points, err := client.Get(ctx, &qdrant.GetPoints{
		CollectionName: name,
		Ids:            []*qdrant.PointId{qdrant.NewIDNum(1)},
		WithPayload:    qdrant.NewWithPayload(true),
	})
	require.NoError(t, err)
	require.Len(t, points, 1)

	result, err = client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: name,
		Wait:           qdrant.PtrOf(true),
		Points:         qdrant.NewPointsSelector(qdrant.NewIDNum(2)),
	})
	require.NoError(t, err)
	require.Equal(t, qdrant.UpdateStatus_Completed, result.GetStatus())
}
