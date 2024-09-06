// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: collections_service.proto

package qdrant

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Collections_Get_FullMethodName                          = "/qdrant.Collections/Get"
	Collections_List_FullMethodName                         = "/qdrant.Collections/List"
	Collections_Create_FullMethodName                       = "/qdrant.Collections/Create"
	Collections_Update_FullMethodName                       = "/qdrant.Collections/Update"
	Collections_Delete_FullMethodName                       = "/qdrant.Collections/Delete"
	Collections_UpdateAliases_FullMethodName                = "/qdrant.Collections/UpdateAliases"
	Collections_ListCollectionAliases_FullMethodName        = "/qdrant.Collections/ListCollectionAliases"
	Collections_ListAliases_FullMethodName                  = "/qdrant.Collections/ListAliases"
	Collections_CollectionClusterInfo_FullMethodName        = "/qdrant.Collections/CollectionClusterInfo"
	Collections_CollectionExists_FullMethodName             = "/qdrant.Collections/CollectionExists"
	Collections_UpdateCollectionClusterSetup_FullMethodName = "/qdrant.Collections/UpdateCollectionClusterSetup"
	Collections_CreateShardKey_FullMethodName               = "/qdrant.Collections/CreateShardKey"
	Collections_DeleteShardKey_FullMethodName               = "/qdrant.Collections/DeleteShardKey"
)

// CollectionsClient is the client API for Collections service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollectionsClient interface {
	// Get detailed information about specified existing collection
	Get(ctx context.Context, in *GetCollectionInfoRequest, opts ...grpc.CallOption) (*GetCollectionInfoResponse, error)
	// Get list name of all existing collections
	List(ctx context.Context, in *ListCollectionsRequest, opts ...grpc.CallOption) (*ListCollectionsResponse, error)
	// Create new collection with given parameters
	Create(ctx context.Context, in *CreateCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error)
	// Update parameters of the existing collection
	Update(ctx context.Context, in *UpdateCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error)
	// Drop collection and all associated data
	Delete(ctx context.Context, in *DeleteCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error)
	// Update Aliases of the existing collection
	UpdateAliases(ctx context.Context, in *ChangeAliases, opts ...grpc.CallOption) (*CollectionOperationResponse, error)
	// Get list of all aliases for a collection
	ListCollectionAliases(ctx context.Context, in *ListCollectionAliasesRequest, opts ...grpc.CallOption) (*ListAliasesResponse, error)
	// Get list of all aliases for all existing collections
	ListAliases(ctx context.Context, in *ListAliasesRequest, opts ...grpc.CallOption) (*ListAliasesResponse, error)
	// Get cluster information for a collection
	CollectionClusterInfo(ctx context.Context, in *CollectionClusterInfoRequest, opts ...grpc.CallOption) (*CollectionClusterInfoResponse, error)
	// Check the existence of a collection
	CollectionExists(ctx context.Context, in *CollectionExistsRequest, opts ...grpc.CallOption) (*CollectionExistsResponse, error)
	// Update cluster setup for a collection
	UpdateCollectionClusterSetup(ctx context.Context, in *UpdateCollectionClusterSetupRequest, opts ...grpc.CallOption) (*UpdateCollectionClusterSetupResponse, error)
	// Create shard key
	CreateShardKey(ctx context.Context, in *CreateShardKeyRequest, opts ...grpc.CallOption) (*CreateShardKeyResponse, error)
	// Delete shard key
	DeleteShardKey(ctx context.Context, in *DeleteShardKeyRequest, opts ...grpc.CallOption) (*DeleteShardKeyResponse, error)
}

type collectionsClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectionsClient(cc grpc.ClientConnInterface) CollectionsClient {
	return &collectionsClient{cc}
}

func (c *collectionsClient) Get(ctx context.Context, in *GetCollectionInfoRequest, opts ...grpc.CallOption) (*GetCollectionInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCollectionInfoResponse)
	err := c.cc.Invoke(ctx, Collections_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) List(ctx context.Context, in *ListCollectionsRequest, opts ...grpc.CallOption) (*ListCollectionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCollectionsResponse)
	err := c.cc.Invoke(ctx, Collections_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) Create(ctx context.Context, in *CreateCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionOperationResponse)
	err := c.cc.Invoke(ctx, Collections_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) Update(ctx context.Context, in *UpdateCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionOperationResponse)
	err := c.cc.Invoke(ctx, Collections_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) Delete(ctx context.Context, in *DeleteCollection, opts ...grpc.CallOption) (*CollectionOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionOperationResponse)
	err := c.cc.Invoke(ctx, Collections_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) UpdateAliases(ctx context.Context, in *ChangeAliases, opts ...grpc.CallOption) (*CollectionOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionOperationResponse)
	err := c.cc.Invoke(ctx, Collections_UpdateAliases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) ListCollectionAliases(ctx context.Context, in *ListCollectionAliasesRequest, opts ...grpc.CallOption) (*ListAliasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAliasesResponse)
	err := c.cc.Invoke(ctx, Collections_ListCollectionAliases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) ListAliases(ctx context.Context, in *ListAliasesRequest, opts ...grpc.CallOption) (*ListAliasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAliasesResponse)
	err := c.cc.Invoke(ctx, Collections_ListAliases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) CollectionClusterInfo(ctx context.Context, in *CollectionClusterInfoRequest, opts ...grpc.CallOption) (*CollectionClusterInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionClusterInfoResponse)
	err := c.cc.Invoke(ctx, Collections_CollectionClusterInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) CollectionExists(ctx context.Context, in *CollectionExistsRequest, opts ...grpc.CallOption) (*CollectionExistsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CollectionExistsResponse)
	err := c.cc.Invoke(ctx, Collections_CollectionExists_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) UpdateCollectionClusterSetup(ctx context.Context, in *UpdateCollectionClusterSetupRequest, opts ...grpc.CallOption) (*UpdateCollectionClusterSetupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCollectionClusterSetupResponse)
	err := c.cc.Invoke(ctx, Collections_UpdateCollectionClusterSetup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) CreateShardKey(ctx context.Context, in *CreateShardKeyRequest, opts ...grpc.CallOption) (*CreateShardKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateShardKeyResponse)
	err := c.cc.Invoke(ctx, Collections_CreateShardKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionsClient) DeleteShardKey(ctx context.Context, in *DeleteShardKeyRequest, opts ...grpc.CallOption) (*DeleteShardKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteShardKeyResponse)
	err := c.cc.Invoke(ctx, Collections_DeleteShardKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

