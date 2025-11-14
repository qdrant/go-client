package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

const (
	collectionName = "test_collection"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

//nolint:funlen,mnd // a long function with magic numbers is acceptable for this example
func run() error {
	// Create new client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost", // Can be omitted, default is "localhost"
		Port: 6334,        // Can be omitted, default is 6334
		// APIKey: "<API_KEY>",
		// UseTLS: true,
		// PoolSize: 3,
		// KeepAliveTime: 10,
		// KeepAliveTimeout: 2,
		// TLSConfig: &tls.Config{...},
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
		return fmt.Errorf("could not get health: %w", err)
	}
	log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())
	// Create collection
	defaultSegmentNumber := uint64(2)
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(4),
			Distance: qdrant.Distance_Dot,
		}),
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create collection: %w", err)
	}
	log.Println("Collection", collectionName, "created")
	// List collections
	collections, err := client.ListCollections(ctx)
	if err != nil {
		return fmt.Errorf("could not list collections: %w", err)
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
		return fmt.Errorf("could not upsert points: %w", err)
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
		return fmt.Errorf("could not retrieve points: %w", err)
	}
	log.Printf("Retrieved points: %s", points)
	// Query the database
	searchedPoints, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
		WithPayload:    qdrant.NewWithPayloadInclude("city"),
	})
	if err != nil {
		return fmt.Errorf("could not search points: %w", err)
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
		return fmt.Errorf("could not search points: %w", err)
	}
	log.Printf("Found points: %s", filteredPoints)
	// Delete collection
	err = client.DeleteCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("could not delete collection: %w", err)
	}
	log.Println("Collection", collectionName, "deleted")
	return nil
}
