package qdrant_test

import (
	"context"
	"log"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

//nolint:testableexamples // there is no qdrant database running
func Example() {
	var (
		collectionName              = "test_collection"
		vectorSize           uint64 = 4
		distance                    = qdrant.Distance_Dot
		defaultSegmentNumber uint64 = 2
	)

	// Create new client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost", // Can be omitted, default is "localhost"
		Port: 6334,        // Can be omitted, default is 6334
		// APIKey: "<API_KEY>",
		// UseTLS: true,
		// TLSConfig: &tls.Config{},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	// Get a context for a minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// Execute health check
	healthCheckResult, err := client.HealthCheck(ctx)
	if err != nil {
		//nolint:gocritic // log.Fatalf is used for simplicity please return an error in real code
		log.Fatalf("Could not get health: %v", err)
	}
	log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())
	// Delete collection
	err = client.DeleteCollection(ctx, collectionName)
	if err != nil {
		log.Fatalf("Could not delete collection: %v", err)
	}
	log.Println("Collection", collectionName, "deleted")
	// Create collection
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: distance,
		}),
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		log.Fatalf("Could not create collection: %v", err)
	}
	log.Println("Collection", collectionName, "created")
	// List collections
	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Fatalf("Could not list collections: %v", err)
	}
	log.Printf("List of collections: %s", &collections)

	// Upsert some data
	waitUpsert := true
	upsertPoints := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
			Payload: qdrant.NewValueMap(map[string]any{
				"city":    "Berlin",
				"country": "Germany",
				"count":   1000000,
				"square":  12.5,
			}),
		},
		{
			Id:      qdrant.NewIDNum(2),
			Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
			Payload: qdrant.NewValueMap(map[string]any{
				"city":    "Berlin",
				"country": "London",
			}),
		},
		{
			Id:      qdrant.NewIDNum(3),
			Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
			Payload: qdrant.NewValueMap(map[string]any{
				"city": []any{"Berlin", "London"},
			}),
		},
		{
			Id:      qdrant.NewID("58384991-3295-4e21-b711-fd3b94fa73e3"),
			Vectors: qdrant.NewVectors(0.35, 0.08, 0.11, 0.44),
			Payload: qdrant.NewValueMap(map[string]any{
				"bool":   true,
				"list":   []any{true, 1, "string"},
				"count":  1000000,
				"square": 12.5,
			}),
		},
	}
	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})
	if err != nil {
		log.Fatalf("Could not upsert points: %v", err)
	}
	log.Println("Upsert", len(upsertPoints), "points")
	// Get points
	points, err := client.Get(ctx, &qdrant.GetPoints{
		CollectionName: collectionName,
		Ids: []*qdrant.PointId{
			qdrant.NewIDNum(1),
			qdrant.NewIDNum(2),
		},
	})
	if err != nil {
		log.Fatalf("Could not retrieve points: %v", err)
	}
	log.Printf("Retrieved points: %s", points)
	// Query the database
	searchedPoints, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
		WithPayload:    qdrant.NewWithPayloadInclude("city"),
	})
	if err != nil {
		log.Fatalf("Could not search points: %v", err)
	}
	log.Printf("Found points: %s", searchedPoints)
	// Query again (with filter)
	filteredPoints, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
		Filter: &qdrant.Filter{
			Should: []*qdrant.Condition{
				qdrant.NewMatchKeyword("city", "Berlin"),
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not search points: %v", err)
	}
	log.Printf("Found points: %s", filteredPoints)
}

//nolint:testableexamples // the example host is invalid
func Example_authentication() {
	// Create new client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   "xyz-example.eu-central.aws.cloud.qdrant.io",
		Port:   6334,
		APIKey: "<paste-your-api-key-here>",
		UseTLS: true,
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		log.Fatalf("could not instantiate: %v", err)
	}
	defer client.Close()
	// List collections
	collections, err := client.ListCollections(context.Background())
	if err != nil {
		//nolint:gocritic // log.Fatalf is used for simplicity please return an error in real code
		log.Fatalf("could not get collections: %v", err)
	}
	log.Printf("List of collections: %v", collections)
}
