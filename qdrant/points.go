package qdrant

import (
	"context"
)

// Perform insert + updates on points. If a point with a given ID already exists - it will be overwritten.
func (c *Client) Upsert(ctx context.Context, request *UpsertPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().Upsert(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Upsert", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Delete points from a collection by IDs or payload filters.
func (c *Client) Delete(ctx context.Context, request *DeletePoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().Delete(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Delete", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Get points from a collection by IDs.
func (c *Client) Get(ctx context.Context, request *GetPoints) ([]*RetrievedPoint, error) {
	resp, err := c.GetPointsClient().Get(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Get", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Iterate over all or filtered points in a collection.
func (c *Client) Scroll(ctx context.Context, request *ScrollPoints) ([]*RetrievedPoint, error) {
	resp, err := c.GetPointsClient().Scroll(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Scroll", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Update vectors for points in a collection.
func (c *Client) UpdateVectors(ctx context.Context, request *UpdatePointVectors) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().UpdateVectors(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "UpdateVectors", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Delete vectors from points in a collection.
func (c *Client) DeleteVectors(ctx context.Context, request *DeletePointVectors) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeleteVectors(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeleteVectors", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Set payload fields for points in a collection.
func (c *Client) SetPayload(ctx context.Context, request *SetPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().SetPayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "SetPayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Overwrite payload for points in a collection.
func (c *Client) OverwritePayload(ctx context.Context, request *SetPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().OverwritePayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "OverwritePayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Delete payload fields from points in a collection.
func (c *Client) DeletePayload(ctx context.Context, request *DeletePayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeletePayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeletePayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Clear payload fields for points in a collection.
func (c *Client) ClearPayload(ctx context.Context, request *ClearPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().ClearPayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "ClearPayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Create index for a payload field.
func (c *Client) CreateFieldIndex(ctx context.Context, request *CreateFieldIndexCollection) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().CreateFieldIndex(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "CreateFieldIndex", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Delete index for a payload field.
func (c *Client) DeleteFieldIndex(ctx context.Context, request *DeleteFieldIndexCollection) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeleteFieldIndex(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeleteFieldIndex", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Count points in collection with given filtering conditions.
// Gets the total count if no filter is provided.
func (c *Client) Count(ctx context.Context, request *CountPoints) (uint64, error) {
	resp, err := c.GetPointsClient().Count(ctx, request)
	if err != nil {
		return 0, newQdrantErr(err, "Count", request.GetCollectionName())
	}
	return resp.GetResult().GetCount(), nil
}

// Perform multiple update operations in one request.
func (c *Client) UpdateBatch(ctx context.Context, request *UpdateBatchPoints) ([]*UpdateResult, error) {
	resp, err := c.GetPointsClient().UpdateBatch(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "UpdateBatch", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Universally query points.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
func (c *Client) Query(ctx context.Context, request *QueryPoints) ([]*ScoredPoint, error) {
	resp, err := c.GetPointsClient().Query(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Query", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Universally query points in a batch.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
func (c *Client) QueryBatch(ctx context.Context, request *QueryBatchPoints) ([]*BatchResult, error) {
	resp, err := c.GetPointsClient().QueryBatch(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "QueryBatch", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Universally query points grouped by a payload field.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
func (c *Client) QueryGroups(ctx context.Context, request *QueryPointGroups) ([]*PointGroup, error) {
	resp, err := c.GetPointsClient().QueryGroups(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "QueryGroups", request.GetCollectionName())
	}
	return resp.GetResult().GetGroups(), nil
}
