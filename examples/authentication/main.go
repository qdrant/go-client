package main

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

//nolint:mnd,gocritic // magic numbers is acceptable for this example, as well as exit with fatal.
func main() {
	// Create new client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   "xyz-example.eu-central.aws.cloud.qdrant.io",
		Port:   6334,
		APIKey: "<paste-your-api-key-here>",
		UseTLS: true,
		Cloud:  true,
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
		log.Fatalf("could not get collections: %v", err)
	}
	log.Printf("List of collections: %v", collections)
}
