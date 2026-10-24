package main

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
	"github.com/qdrant/go-client/qdrant/serverless"
)

//nolint:gocritic // exit with fatal is acceptable for this example.
func main() {
	ctx := context.Background()

	client, err := serverless.NewClient(&qdrant.Config{
		Host:   "serverless.example.cloud.qdrant.io",
		APIKey: "<paste-your-api-key-here>",
	})
	if err != nil {
		log.Fatalf("could not instantiate: %v", err)
	}
	defer client.Close()

	collectionName := "my-collection"
	createCollection(ctx, client, collectionName)
	upsertAndQuery(ctx, client, collectionName)

	if _, err := client.DeleteCollection(ctx, collectionName); err != nil {
		log.Fatalf("could not delete collection: %v", err)
	}
}

//nolint:mnd // magic numbers is acceptable for this example.
func createCollection(ctx context.Context, client *serverless.Client, collectionName string) {
	exists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		log.Fatalf("could not check collection: %v", err)
	}
	if exists {
		if _, err := client.DeleteCollection(ctx, collectionName); err != nil {
			log.Fatalf("could not delete collection: %v", err)
		}
	}

	created, err := client.CreateCollection(ctx, collectionName, &serverless.CollectionConfig{
		DenseVectors: map[string]*serverless.DenseVectorConfig{
			"": {Size: 4, Distance: serverless.Distance_COSINE},
		},
		PayloadIndexes: map[string]*serverless.PayloadIndexConfig{
			"color": serverless.NewPayloadIndexConfigKeyword(nil),
		},
	})
	if err != nil {
		log.Fatalf("could not create collection: %v", err)
	}
	log.Printf("Create collection: %v", created)

	collections, err := client.ListCollections(ctx, nil)
	if err != nil {
		log.Fatalf("could not list collections: %v", err)
	}
	log.Printf("List of collections: %v", collections.GetCollections())

	info, err := client.GetCollection(ctx, collectionName)
	if err != nil {
		log.Fatalf("could not get collection: %v", err)
	}
	log.Printf("Collection info: %v", info)
}

//nolint:mnd // magic numbers is acceptable for this example.
func upsertAndQuery(ctx context.Context, client *serverless.Client, collectionName string) {
	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
				Payload: qdrant.NewValueMap(map[string]any{"color": "red"}),
			},
			{
				Id:      qdrant.NewIDNum(2),
				Vectors: qdrant.NewVectors(0.4, 0.3, 0.2, 0.1),
				Payload: qdrant.NewValueMap(map[string]any{"color": "blue"}),
			},
		},
	})
	if err != nil {
		log.Fatalf("could not upsert points: %v", err)
	}

	points, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(0.1, 0.2, 0.3, 0.4),
		Limit:          qdrant.PtrOf(uint64(10)),
	})
	if err != nil {
		log.Fatalf("could not query points: %v", err)
	}
	log.Printf("Query results: %v", points)
}
