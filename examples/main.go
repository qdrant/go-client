package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr                  = flag.String("addr", "localhost:6334", "the address to connect to")
	collectionName        = "test_collection"
	vectorSize     uint64 = 4
	distance              = qdrant.Distance_Dot
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), *addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// create grpc collection client
	collections_client := qdrant.NewCollectionsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Check Qdrant version
	qdrantClient := qdrant.NewQdrantClient(conn)
	healthCheckResult, err := qdrantClient.HealthCheck(ctx, &qdrant.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("Could not get health: %v", err)
	} else {
		log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())
	}

	// Delete collection
	_, err = collections_client.Delete(ctx, &qdrant.DeleteCollection{
		CollectionName: collectionName,
	})
	if err != nil {
		log.Fatalf("Could not delete collection: %v", err)
	} else {
		log.Println("Collection", collectionName, "deleted")
	}

	// Create new collection
	var defaultSegmentNumber uint64 = 2
	_, err = collections_client.Create(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &qdrant.VectorsConfig{Config: &qdrant.VectorsConfig_Params{
			Params: &qdrant.VectorParams{
				Size:     vectorSize,
				Distance: distance,
			},
		}},
		OptimizersConfig: &qdrant.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		log.Fatalf("Could not create collection: %v", err)
	} else {
		log.Println("Collection", collectionName, "created")
	}

	// List all created collections
	r, err := collections_client.List(ctx, &qdrant.ListCollectionsRequest{})
	if err != nil {
		log.Fatalf("Could not get collections: %v", err)
	} else {
		log.Printf("List of collections: %s", r.GetCollections())
	}

	// Create points grpc client
	pointsClient := qdrant.NewPointsClient(conn)

	// Create keyword field index
	fieldIndex1Type := qdrant.FieldType_FieldTypeKeyword
	fieldIndex1Name := "city"
	_, err = pointsClient.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: collectionName,
		FieldName:      fieldIndex1Name,
		FieldType:      &fieldIndex1Type,
	})
	if err != nil {
		log.Fatalf("Could not create field index: %v", err)
	} else {
		log.Println("Field index for field", fieldIndex1Name, "created")
	}

	// Create integer field index
	fieldIndex2Type := qdrant.FieldType_FieldTypeInteger
	fieldIndex2Name := "count"
	_, err = pointsClient.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: collectionName,
		FieldName:      fieldIndex2Name,
		FieldType:      &fieldIndex2Type,
	})
	if err != nil {
		log.Fatalf("Could not create field index: %v", err)
	} else {
		log.Println("Field index for field", fieldIndex2Name, "created")
	}

	// Upsert points
	waitUpsert := true
	upsertPoints := []*qdrant.PointStruct{
		{
			// Point Id is number or UUID
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 1},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.05, 0.61, 0.76, 0.74}}}},
			Payload: map[string]*qdrant.Value{
				"city": {
					Kind: &qdrant.Value_StringValue{StringValue: "Berlin"},
				},
				"country": {
					Kind: &qdrant.Value_StringValue{StringValue: "Germany"},
				},
				"count": {
					Kind: &qdrant.Value_IntegerValue{IntegerValue: 1000000},
				},
				"square": {
					Kind: &qdrant.Value_DoubleValue{DoubleValue: 12.5},
				},
			},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 2},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.19, 0.81, 0.75, 0.11}}}},
			Payload: map[string]*qdrant.Value{
				"city": {
					Kind: &qdrant.Value_ListValue{
						ListValue: &qdrant.ListValue{
							Values: []*qdrant.Value{
								{
									Kind: &qdrant.Value_StringValue{StringValue: "Berlin"},
								},
								{
									Kind: &qdrant.Value_StringValue{StringValue: "London"},
								},
							},
						},
					},
				},
			},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 3},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.36, 0.55, 0.47, 0.94}}}},
			Payload: map[string]*qdrant.Value{
				"city": {
					Kind: &qdrant.Value_ListValue{
						ListValue: &qdrant.ListValue{
							Values: []*qdrant.Value{
								{
									Kind: &qdrant.Value_StringValue{StringValue: "Berlin"},
								},
								{
									Kind: &qdrant.Value_StringValue{StringValue: "Moscow"},
								},
							},
						},
					},
				},
			},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 4},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.18, 0.01, 0.85, 0.80}}}},
			Payload: map[string]*qdrant.Value{
				"city": {
					Kind: &qdrant.Value_ListValue{
						ListValue: &qdrant.ListValue{
							Values: []*qdrant.Value{
								{
									Kind: &qdrant.Value_StringValue{StringValue: "London"},
								},
								{
									Kind: &qdrant.Value_StringValue{StringValue: "Moscow"},
								},
							},
						},
					},
				},
			},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 5},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.24, 0.18, 0.22, 0.44}}}},
			Payload: map[string]*qdrant.Value{
				"count": {
					Kind: &qdrant.Value_ListValue{
						ListValue: &qdrant.ListValue{
							Values: []*qdrant.Value{
								{
									Kind: &qdrant.Value_IntegerValue{IntegerValue: 0},
								},
							},
						},
					},
				},
			},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{Num: 6},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
			Payload: map[string]*qdrant.Value{},
		},
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: "58384991-3295-4e21-b711-fd3b94fa73e3"},
			},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
			Payload: map[string]*qdrant.Value{},
		},
	}
	_, err = pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})
	if err != nil {
		log.Fatalf("Could not upsert points: %v", err)
	} else {
		log.Println("Upsert", len(upsertPoints), "points")
	}

	// Retrieve points by ids
	pointsById, err := pointsClient.Get(ctx, &qdrant.GetPoints{
		CollectionName: collectionName,
		Ids: []*qdrant.PointId{
			{PointIdOptions: &qdrant.PointId_Num{Num: 1}},
			{PointIdOptions: &qdrant.PointId_Num{Num: 2}},
		},
	})
	if err != nil {
		log.Fatalf("Could not retrieve points: %v", err)
	} else {
		log.Printf("Retrieved points: %s", pointsById.GetResult())
	}

	// Unfiltered search
	unfilteredSearchResult, err := pointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: collectionName,
		Vector:         []float32{0.2, 0.1, 0.9, 0.7},
		Limit:          3,
		// Include all payload and vectors in the search result
		WithVectors: &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: true}},
		WithPayload: &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		log.Fatalf("Could not search points: %v", err)
	} else {
		log.Printf("Found points: %s", unfilteredSearchResult.GetResult())
	}

	// filtered search
	filteredSearchResult, err := pointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: collectionName,
		Vector:         []float32{0.2, 0.1, 0.9, 0.7},
		Limit:          3,
		Filter: &qdrant.Filter{
			Should: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: "city",
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Keyword{
									Keyword: "London",
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not search points: %v", err)
	} else {
		log.Printf("Found points: %s", filteredSearchResult.GetResult())
	}
}
