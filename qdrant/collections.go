package qdrant

import (
	"context"
	"errors"
)

// Check the existence of a collection.
func (c *Client) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	resp, err := c.GetCollectionsClient().CollectionExists(ctx, &CollectionExistsRequest{
		CollectionName: collectionName,
	})
	if err != nil {
		return false, newQdrantErr(err, "CollectionExists", collectionName)
	}
	return resp.GetResult().GetExists(), nil
}

// Get detailed information about specified existing collection.
func (c *Client) GetCollection(ctx context.Context, collectionName string) (*CollectionInfo, error) {
	resp, err := c.GetCollectionsClient().Get(ctx, &GetCollectionInfoRequest{
		CollectionName: collectionName,
	})
	if err != nil {
		return nil, newQdrantErr(err, "GetCollection", collectionName)
	}
	return resp.GetResult(), nil
}

// Get names of all existing collections.
func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	resp, err := c.GetCollectionsClient().List(ctx, &ListCollectionsRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "ListCollections")
	}
	var collections []string
	for _, collection := range resp.GetCollections() {
		collections = append(collections, collection.GetName())
	}
	return collections, nil
}

// Create new collection with given parameters.
func (c *Client) CreateCollection(ctx context.Context, request *CreateCollection) error {
	_, err := c.GetCollectionsClient().Create(ctx, request)
	if err != nil {
		return newQdrantErr(err, "CreateCollection", request.GetCollectionName())
	}
	return nil
}

// Update parameters of the existing collection.
func (c *Client) UpdateCollection(ctx context.Context, request *UpdateCollection) error {
	_, err := c.GetCollectionsClient().Update(ctx, request)
	if err != nil {
		return newQdrantErr(err, "UpdateCollection", request.GetCollectionName())
	}
	return nil
}

// Drop collection and all associated data.
func (c *Client) DeleteCollection(ctx context.Context, collectionName string) error {
	res, err := c.GetCollectionsClient().Delete(ctx, &DeleteCollection{
		CollectionName: collectionName,
	})

	if err != nil {
		return newQdrantErr(err, "DeleteCollection", collectionName)
	}

	if !res.GetResult() {
		return newQdrantErr(errors.New("failed to delete collection"), "DeleteCollection", collectionName)
	}

	return nil
}

// Create an alias for a collection.
func (c *Client) CreateAlias(ctx context.Context, aliasName, collectionName string) error {
	_, err := c.GetCollectionsClient().UpdateAliases(ctx, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_CreateAlias{
					CreateAlias: &CreateAlias{
						CollectionName: collectionName,
						AliasName:      aliasName,
					},
				},
			},
		},
	})
	if err != nil {
		return newQdrantErr(err, "CreateAlias", collectionName)
	}
	return nil
}

// Delete an alias.
func (c *Client) DeleteAlias(ctx context.Context, aliasName string) error {
	_, err := c.GetCollectionsClient().UpdateAliases(ctx, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_DeleteAlias{
					DeleteAlias: &DeleteAlias{
						AliasName: aliasName,
					},
				},
			},
		},
	})
	if err != nil {
		return newQdrantErr(err, "DeleteAlias", aliasName)
	}
	return nil
}

// Rename an alias.
func (c *Client) RenameAlias(ctx context.Context, oldAliasName, newAliasName string) error {
	_, err := c.GetCollectionsClient().UpdateAliases(ctx, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_RenameAlias{
					RenameAlias: &RenameAlias{
						OldAliasName: oldAliasName,
						NewAliasName: newAliasName,
					},
				},
			},
		},
	})
	if err != nil {
		return newQdrantErr(err, "RenameAlias")
	}
	return nil
}

// List all aliases for a collection.
func (c *Client) ListCollectionAliases(ctx context.Context, collectionName string) ([]string, error) {
	resp, err := c.GetCollectionsClient().ListCollectionAliases(ctx, &ListCollectionAliasesRequest{
		CollectionName: collectionName,
	})
	if err != nil {
		return nil, newQdrantErr(err, "ListCollectionAliases", collectionName)
	}
	var aliases []string
	for _, alias := range resp.GetAliases() {
		aliases = append(aliases, alias.GetAliasName())
	}
	return aliases, nil
}

// List all aliases.
func (c *Client) ListAliases(ctx context.Context) ([]*AliasDescription, error) {
	resp, err := c.GetCollectionsClient().ListAliases(ctx, &ListAliasesRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "ListAliases")
	}
	return resp.GetAliases(), nil
}

// Update aliases.
func (c *Client) UpdateAliases(ctx context.Context, actions []*AliasOperations) error {
	_, err := c.GetCollectionsClient().UpdateAliases(ctx, &ChangeAliases{
		Actions: actions,
	})
	if err != nil {
		return newQdrantErr(err, "UpdateAliases")
	}
	return nil
}

// Create a shard key for a collection.
func (c *Client) CreateShardKey(ctx context.Context, collectionName string, request *CreateShardKey) error {
	_, err := c.GetCollectionsClient().CreateShardKey(ctx, &CreateShardKeyRequest{
		CollectionName: collectionName,
		Request:        request,
	})
	if err != nil {
		return newQdrantErr(err, "CreateShardKey", collectionName)
	}
	return nil
}

// Delete a shard key for a collection.
func (c *Client) DeleteShardKey(ctx context.Context, collectionName string, request *DeleteShardKey) error {
	_, err := c.GetCollectionsClient().DeleteShardKey(ctx, &DeleteShardKeyRequest{
		CollectionName: collectionName,
		Request:        request,
	})
	if err != nil {
		return newQdrantErr(err, "DeleteShardKey", collectionName)
	}
	return nil
}
