<p align="center">
  <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://github.com/qdrant/qdrant/raw/master/docs/logo-dark.svg">
      <source media="(prefers-color-scheme: light)" srcset="https://github.com/qdrant/qdrant/raw/master/docs/logo-light.svg">
      <img height="100" alt="Qdrant" src="https://github.com/qdrant/qdrant/raw/master/docs/logo.svg">
  </picture>
  <img height="150" width="100" src="./examples/512px-Go_Logo_Blue.svg.png" alt="Golang">
</p>

<p align="center">
    <b>Go client for the <a href="https://github.com/qdrant/qdrant">Qdrant</a> vector search engine.</b>
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/qdrant/go-client"><img src="https://img.shields.io/badge/Docs-godoc-success" alt="Godoc"></a>
    <a href="https://github.com/qdrant/go-client/actions/workflows/ci.yml"><img src="https://github.com/qdrant/go-client/actions/workflows/ci.yml/badge.svg?branch=master" alt="Tests"></a>
    <a href="https://github.com/qdrant/go-client/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-success" alt="Apache 2.0 License"></a>
    <a href="https://qdrant.to/discord"><img src="https://img.shields.io/badge/Discord-Qdrant-5865F2.svg?logo=discord" alt="Discord"></a>
    <a href="https://qdrant.to/roadmap"><img src="https://img.shields.io/badge/Roadmap-2025-bc1439.svg" alt="Roadmap 2025"></a>
</p>

Go client library with handy utilities for interfacing with [Qdrant](https://qdrant.tech/).

## 📥 Installation

```bash
go get -u github.com/qdrant/go-client
```

## 📖 Documentation

- Usage examples are available throughout the [Qdrant documentation](https://qdrant.tech/documentation/quick-start/) and [API Reference](https://api.qdrant.tech/).
- [Godoc Reference](https://pkg.go.dev/github.com/qdrant/go-client)

## 🔌 Getting started

### Creating a client

A client can be instantiated with

```go
import "github.com/qdrant/go-client/qdrant"

client, err := qdrant.NewClient(&qdrant.Config{
  Host: "localhost",
  Port: 6334,
})
```

Which creates a client that will connect to Qdrant on <http://localhost:6334>.

Internally, the high-level client uses a low-level gRPC client to interact with
Qdrant. `qdrant.Config` provides additional options to control how the gRPC
client is configured. The following example configures API key authentication with TLS:

```go
import "github.com/qdrant/go-client/qdrant"

client, err := qdrant.NewClient(&qdrant.Config{
	Host:   "xyz-example.eu-central.aws.cloud.qdrant.io",
	Port:   6334,
	APIKey: "<your-api-key>",
	UseTLS: true,  // uses default config with minimum TLS version set to 1.3
	// TLSConfig: &tls.Config{...},
	// GrpcOptions: []grpc.DialOption{},
})
```

### Working with collections

Once a client has been created, create a new collection

```go
import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

client.CreateCollection(context.Background(), &qdrant.CreateCollection{
	CollectionName: "{collection_name}",
	VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
		Size:     4,
		Distance: qdrant.Distance_Cosine,
	}),
})
```

Insert vectors into the collection

```go
operationInfo, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
	CollectionName: "{collection_name}",
	Points: []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
			Payload: qdrant.NewValueMap(map[string]any{"city": "London"}),
		},
		{
			Id:      qdrant.NewIDNum(2),
			Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
			Payload: qdrant.NewValueMap(map[string]any{"age": 32}),
		},
		{
			Id:      qdrant.NewIDNum(3),
			Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
			Payload: qdrant.NewValueMap(map[string]any{"vegan": true}),
		},
	},
})
```

Search for similar vectors

```go
searchResult, err := client.Query(context.Background(), &qdrant.QueryPoints{
	CollectionName: "{collection_name}",
	Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
})
```

Search for similar vectors with a filtering condition

```go
searchResult, err := client.Query(context.Background(), &qdrant.QueryPoints{
	CollectionName: "{collection_name}",
	Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
	Filter: &qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewMatch("city", "London"),
		},
	},
	WithPayload: qdrant.NewWithPayload(true),
})
```

## ⚖️ LICENSE

[Apache 2.0](https://github.com/qdrant/go-client/blob/master/LICENSE)
