package go_client_test

import (
	context "context"
	"log"
	"testing"
	"time"

	pb "github.com/qdrant/go-client/qdrant"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/qdrant"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewQdrantClient(t *testing.T) {
	var (
		collectionName        = "test_collection"
		vectorSize     uint64 = 4
		distance              = pb.Distance_Dot
	)

	c, err := qdrant.RunContainer(context.Background(), testcontainers.WithImage("qdrant/qdrant:v1.9.0"))
	if err != nil {
		t.Fatalf("Could not start qdrant container: %v", err)
	}

	addr, err := c.GRPCEndpoint(context.Background())
	if err != nil {
		t.Fatalf("Could not get qdrant container grpc endpoint: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// create grpc collection client
	collections_client := pb.NewCollectionsClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	qdrantClient := pb.NewQdrantClient(conn)

	t.Run("Check Qdrant version", func(t *testing.T) {
		healthCheckResult, err := qdrantClient.HealthCheck(ctx, &pb.HealthCheckRequest{})
		if err != nil {
			t.Fatalf("Could not get health: %v", err)
		} else {
			log.Printf("Qdrant version: %s", healthCheckResult.GetVersion())
		}
	})

	t.Run("Delete collection", func(t *testing.T) {
		_, err = collections_client.Delete(ctx, &pb.DeleteCollection{
			CollectionName: collectionName,
		})
		if err != nil {
			t.Fatalf("Could not delete collection: %v", err)
		} else {
			log.Println("Collection", collectionName, "deleted")
		}
	})

	t.Run("Create new collection", func(t *testing.T) {
		var defaultSegmentNumber uint64 = 2
		_, err = collections_client.Create(ctx, &pb.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     vectorSize,
					Distance: distance,
				},
			}},
			OptimizersConfig: &pb.OptimizersConfigDiff{
				DefaultSegmentNumber: &defaultSegmentNumber,
			},
		})
		if err != nil {
			t.Fatalf("Could not create collection: %v", err)
		} else {
			log.Println("Collection", collectionName, "created")
		}
	})

	t.Run("List all created collections", func(t *testing.T) {
		r, err := collections_client.List(ctx, &pb.ListCollectionsRequest{})
		if err != nil {
			t.Fatalf("Could not get collections: %v", err)
		} else {
			log.Printf("List of collections: %s", r.GetCollections())
		}
	})

	// Create points grpc client
	pointsClient := pb.NewPointsClient(conn)

	t.Run("Create keyword field index", func(t *testing.T) {
		fieldIndex1Type := pb.FieldType_FieldTypeKeyword
		fieldIndex1Name := "city"
		_, err = pointsClient.CreateFieldIndex(ctx, &pb.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      fieldIndex1Name,
			FieldType:      &fieldIndex1Type,
		})
		if err != nil {
			t.Fatalf("Could not create field index: %v", err)
		} else {
			log.Println("Field index for field", fieldIndex1Name, "created")
		}
	})

	t.Run("Create integer field index", func(t *testing.T) {
		fieldIndex2Type := pb.FieldType_FieldTypeInteger
		fieldIndex2Name := "count"
		_, err = pointsClient.CreateFieldIndex(ctx, &pb.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      fieldIndex2Name,
			FieldType:      &fieldIndex2Type,
		})
		if err != nil {
			t.Fatalf("Could not create field index: %v", err)
		} else {
			log.Println("Field index for field", fieldIndex2Name, "created")
		}
	})

	t.Run("Upsert points", func(t *testing.T) {
		waitUpsert := true
		upsertPoints := []*pb.PointStruct{
			{
				// Point Id is number or UUID
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 1},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.05, 0.61, 0.76, 0.74}}}},
				Payload: map[string]*pb.Value{
					"city": {
						Kind: &pb.Value_StringValue{StringValue: "Berlin"},
					},
					"country": {
						Kind: &pb.Value_StringValue{StringValue: "Germany"},
					},
					"count": {
						Kind: &pb.Value_IntegerValue{IntegerValue: 1000000},
					},
					"square": {
						Kind: &pb.Value_DoubleValue{DoubleValue: 12.5},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 2},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.19, 0.81, 0.75, 0.11}}}},
				Payload: map[string]*pb.Value{
					"city": {
						Kind: &pb.Value_ListValue{
							ListValue: &pb.ListValue{
								Values: []*pb.Value{
									{
										Kind: &pb.Value_StringValue{StringValue: "Berlin"},
									},
									{
										Kind: &pb.Value_StringValue{StringValue: "London"},
									},
								},
							},
						},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 3},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.36, 0.55, 0.47, 0.94}}}},
				Payload: map[string]*pb.Value{
					"city": {
						Kind: &pb.Value_ListValue{
							ListValue: &pb.ListValue{
								Values: []*pb.Value{
									{
										Kind: &pb.Value_StringValue{StringValue: "Berlin"},
									},
									{
										Kind: &pb.Value_StringValue{StringValue: "Moscow"},
									},
								},
							},
						},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 4},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.18, 0.01, 0.85, 0.80}}}},
				Payload: map[string]*pb.Value{
					"city": {
						Kind: &pb.Value_ListValue{
							ListValue: &pb.ListValue{
								Values: []*pb.Value{
									{
										Kind: &pb.Value_StringValue{StringValue: "London"},
									},
									{
										Kind: &pb.Value_StringValue{StringValue: "Moscow"},
									},
								},
							},
						},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 5},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.24, 0.18, 0.22, 0.44}}}},
				Payload: map[string]*pb.Value{
					"count": {
						Kind: &pb.Value_ListValue{
							ListValue: &pb.ListValue{
								Values: []*pb.Value{
									{
										Kind: &pb.Value_IntegerValue{IntegerValue: 0},
									},
								},
							},
						},
					},
				},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Num{Num: 6},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
				Payload: map[string]*pb.Value{},
			},
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{Uuid: "58384991-3295-4e21-b711-fd3b94fa73e3"},
				},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: []float32{0.35, 0.08, 0.11, 0.44}}}},
				Payload: map[string]*pb.Value{},
			},
		}
		_, err = pointsClient.Upsert(ctx, &pb.UpsertPoints{
			CollectionName: collectionName,
			Wait:           &waitUpsert,
			Points:         upsertPoints,
		})
		if err != nil {
			t.Fatalf("Could not upsert points: %v", err)
		} else {
			log.Println("Upsert", len(upsertPoints), "points")
		}

		// Retrieve points by ids
		pointsById, err := pointsClient.Get(ctx, &pb.GetPoints{
			CollectionName: collectionName,
			Ids: []*pb.PointId{
				{PointIdOptions: &pb.PointId_Num{Num: 1}},
				{PointIdOptions: &pb.PointId_Num{Num: 2}},
			},
		})
		if err != nil {
			t.Fatalf("Could not retrieve points: %v", err)
		} else {
			log.Printf("Retrieved points: %s", pointsById.GetResult())
		}
	})

	t.Run("Unfiltered search", func(t *testing.T) {
		unfilteredSearchResult, err := pointsClient.Search(ctx, &pb.SearchPoints{
			CollectionName: collectionName,
			Vector:         []float32{0.2, 0.1, 0.9, 0.7},
			Limit:          3,
			// Include all payload and vectors in the search result
			WithVectors: &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: true}},
			WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
		})
		if err != nil {
			t.Fatalf("Could not search points: %v", err)
		} else {
			log.Printf("Found points: %s", unfilteredSearchResult.GetResult())
		}
	})

	t.Run("Filtered search", func(t *testing.T) {
		filteredSearchResult, err := pointsClient.Search(ctx, &pb.SearchPoints{
			CollectionName: collectionName,
			Vector:         []float32{0.2, 0.1, 0.9, 0.7},
			Limit:          3,
			Filter: &pb.Filter{
				Should: []*pb.Condition{
					{
						ConditionOneOf: &pb.Condition_Field{
							Field: &pb.FieldCondition{
								Key: "city",
								Match: &pb.Match{
									MatchValue: &pb.Match_Keyword{
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
			t.Fatalf("Could not search points: %v", err)
		} else {
			log.Printf("Found points: %s", filteredSearchResult.GetResult())
		}
	})
}
