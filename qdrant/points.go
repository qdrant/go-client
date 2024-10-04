package qdrant

import (
	"context"
)

// Performs insert + updates on points. If a point with a given ID already exists, it will be overwritten.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The UpsertPoints request containing the points to upsert.
//
// Returns:
//   - *UpdateResult: The result of the upsert operation.
//   - error: An error if the operation fails.
func (c *Client) Upsert(ctx context.Context, request *UpsertPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().Upsert(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Upsert", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Removes points from a collection by IDs or payload filters.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The DeletePoints request specifying which points to delete.
//
// Returns:
//   - *UpdateResult: The result of the delete operation.
//   - error: An error if the operation fails.
func (c *Client) Delete(ctx context.Context, request *DeletePoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().Delete(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Delete", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Retrieves points from a collection by IDs.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The GetPoints request specifying which points to retrieve.
//
// Returns:
//   - []*RetrievedPoint: A slice of retrieved points.
//   - error: An error if the operation fails.
func (c *Client) Get(ctx context.Context, request *GetPoints) ([]*RetrievedPoint, error) {
	resp, err := c.GetPointsClient().Get(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Get", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Iterates over all or filtered points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The ScrollPoints request specifying the scroll parameters.
//
// Returns:
//   - []*RetrievedPoint: A slice of retrieved points.
//   - error: An error if the operation fails.
func (c *Client) Scroll(ctx context.Context, request *ScrollPoints) ([]*RetrievedPoint, error) {
	resp, err := c.GetPointsClient().Scroll(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Scroll", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Updates vectors for points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The UpdatePointVectors request containing the vectors to update.
//
// Returns:
//   - *UpdateResult: The result of the update operation.
//   - error: An error if the operation fails.
func (c *Client) UpdateVectors(ctx context.Context, request *UpdatePointVectors) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().UpdateVectors(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "UpdateVectors", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Removes vectors from points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The DeletePointVectors request specifying which vectors to delete.
//
// Returns:
//   - *UpdateResult: The result of the delete operation.
//   - error: An error if the operation fails.
func (c *Client) DeleteVectors(ctx context.Context, request *DeletePointVectors) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeleteVectors(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeleteVectors", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Sets payload fields for points in a collection.
// Can be used to add new payload fields or update existing ones.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The SetPayloadPoints request containing the payload to set.
//
// Returns:
//   - *UpdateResult: The result of the set operation.
//   - error: An error if the operation fails.
func (c *Client) SetPayload(ctx context.Context, request *SetPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().SetPayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "SetPayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Overwrites the entire payload for points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The SetPayloadPoints request containing the payload to overwrite.
//
// Returns:
//   - *UpdateResult: The result of the overwrite operation.
//   - error: An error if the operation fails.
func (c *Client) OverwritePayload(ctx context.Context, request *SetPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().OverwritePayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "OverwritePayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Removes payload fields from points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The DeletePayloadPoints request specifying which payload fields to delete.
//
// Returns:
//   - *UpdateResult: The result of the delete operation.
//   - error: An error if the operation fails.
func (c *Client) DeletePayload(ctx context.Context, request *DeletePayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeletePayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeletePayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Removes all payload fields from points in a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The ClearPayloadPoints request specifying which points to clear.
//
// Returns:
//   - *UpdateResult: The result of the clear operation.
//   - error: An error if the operation fails.
func (c *Client) ClearPayload(ctx context.Context, request *ClearPayloadPoints) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().ClearPayload(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "ClearPayload", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Creates an index for a payload field.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The CreateFieldIndexCollection request specifying the field to index.
//
// Returns:
//   - *UpdateResult: The result of the index creation operation.
//   - error: An error if the operation fails.
func (c *Client) CreateFieldIndex(ctx context.Context, request *CreateFieldIndexCollection) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().CreateFieldIndex(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "CreateFieldIndex", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Removes an index for a payload field.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The DeleteFieldIndexCollection request specifying the field index to delete.
//
// Returns:
//   - *UpdateResult: The result of the index deletion operation.
//   - error: An error if the operation fails.
func (c *Client) DeleteFieldIndex(ctx context.Context, request *DeleteFieldIndexCollection) (*UpdateResult, error) {
	resp, err := c.GetPointsClient().DeleteFieldIndex(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "DeleteFieldIndex", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Returns the number of points in a collection with given filtering conditions.
// Gets the total count if no filter is provided.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The CountPoints request containing optional filtering conditions.
//
// Returns:
//   - uint64: The count of points matching the conditions.
//   - error: An error if the operation fails.
func (c *Client) Count(ctx context.Context, request *CountPoints) (uint64, error) {
	resp, err := c.GetPointsClient().Count(ctx, request)
	if err != nil {
		return 0, newQdrantErr(err, "Count", request.GetCollectionName())
	}
	return resp.GetResult().GetCount(), nil
}

// Performs multiple update operations in one request.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The UpdateBatchPoints request containing multiple update operations.
//
// Returns:
//   - []*UpdateResult: A slice of results for each update operation.
//   - error: An error if the operation fails.
func (c *Client) UpdateBatch(ctx context.Context, request *UpdateBatchPoints) ([]*UpdateResult, error) {
	resp, err := c.GetPointsClient().UpdateBatch(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "UpdateBatch", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Performs a universal query on points.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The QueryPoints request containing the query parameters.
//
// Returns:
//   - []*ScoredPoint: A slice of scored points matching the query.
//   - error: An error if the operation fails.
func (c *Client) Query(ctx context.Context, request *QueryPoints) ([]*ScoredPoint, error) {
	resp, err := c.GetPointsClient().Query(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Query", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Performs multiple universal queries on points in a batch.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The QueryBatchPoints request containing multiple query parameters.
//
// Returns:
//   - []*BatchResult: A slice of batch results for each query.
//   - error: An error if the operation fails.
func (c *Client) QueryBatch(ctx context.Context, request *QueryBatchPoints) ([]*BatchResult, error) {
	resp, err := c.GetPointsClient().QueryBatch(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "QueryBatch", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Performs a universal query on points grouped by a payload field.
// Covers all capabilities of search, recommend, discover, filters.
// Also enables hybrid and multi-stage queries.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The QueryPointGroups request containing the query parameters.
//
// Returns:
//   - []*PointGroup: A slice of point groups matching the query.
//   - error: An error if the operation fails.
func (c *Client) QueryGroups(ctx context.Context, request *QueryPointGroups) ([]*PointGroup, error) {
	resp, err := c.GetPointsClient().QueryGroups(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "QueryGroups", request.GetCollectionName())
	}
	return resp.GetResult().GetGroups(), nil
}

// Facets the points that would match a filter, with respect to a payload field.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The FacetCounts request containing the facet parameters.
//
// Returns:
//   - []*FacetHit: A slice of facet hits matching the query. Each hit contains the value and the count for this value
func (c *Client) Facet(ctx context.Context, request *FacetCounts) ([]*FacetHit, error) {
	resp, err := c.GetPointsClient().Facet(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "Facet", request.GetCollectionName())
	}
	return resp.GetHits(), nil
}

// Calculates the distances between a random sample of points.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The SearchMatrixPoints request containing the sample size and other parameters.
//
// Returns:
//   - *SearchMatrixPairs: Pairwise representation of distances.
func (c *Client) SearchMatrixPairs(ctx context.Context, request *SearchMatrixPoints) (*SearchMatrixPairs, error) {
	resp, err := c.GetPointsClient().SearchMatrixPairs(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "SearchMatrixPairs", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}

// Calculates the distances between a random sample of points.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The SearchMatrixPoints request containing the sample size and other parameters.
//
// Returns:
//   - *SearchMatrixOffsets: Parallel lists of offsets within an id list, and the distance between them.
func (c *Client) SearchMatrixOffsets(ctx context.Context, request *SearchMatrixPoints) (*SearchMatrixOffsets, error) {
	resp, err := c.GetPointsClient().SearchMatrixOffsets(ctx, request)
	if err != nil {
		return nil, newQdrantErr(err, "SearchMatrixOffsets", request.GetCollectionName())
	}
	return resp.GetResult(), nil
}
