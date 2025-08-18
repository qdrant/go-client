package main

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

const (
	collectionName = "test_collection"
)

func main() {
	// Create context
	ctx := context.Background()

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
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(128),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		log.Printf("Could not create collection: %v", err)
		return
	}

	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Printf("Could not list collections: %v", err)
		return
	}
	log.Printf("List of collections: %+v", collections)
}
