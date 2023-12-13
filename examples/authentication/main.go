package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "secure.cloud.qdrant.io:6334", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	config := &tls.Config{}
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	collections_client := qdrant.NewCollectionsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	md := metadata.New(map[string]string{"api-key": "secret-key-*******"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	defer cancel()
	r, err := collections_client.List(ctx, &qdrant.ListCollectionsRequest{})
	if err != nil {
		log.Fatalf("could not get collections: %v", err)
	}
	log.Printf("List of collections: %s", r.GetCollections())
}
