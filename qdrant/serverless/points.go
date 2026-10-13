package serverless

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

// Queries points.
func (c *Client) Query(ctx context.Context, request *qdrant.QueryPoints) ([]*qdrant.ScoredPoint, error) {
	return c.points.Query(ctx, request)
}

// Runs several queries in one request.
func (c *Client) QueryBatch(ctx context.Context, request *qdrant.QueryBatchPoints) ([]*qdrant.BatchResult, error) {
	return c.points.QueryBatch(ctx, request)
}

// Queries points grouped by a payload field.
func (c *Client) QueryGroups(ctx context.Context, request *qdrant.QueryPointGroups) ([]*qdrant.PointGroup, error) {
	return c.points.QueryGroups(ctx, request)
}

// Retrieves points by IDs.
func (c *Client) Get(ctx context.Context, request *qdrant.GetPoints) ([]*qdrant.RetrievedPoint, error) {
	return c.points.Get(ctx, request)
}

// Returns a page of points.
func (c *Client) Scroll(ctx context.Context, request *qdrant.ScrollPoints) ([]*qdrant.RetrievedPoint, error) {
	return c.points.Scroll(ctx, request)
}

// Returns a page of points and the offset of the next page.
func (c *Client) ScrollAndOffset(ctx context.Context,
	request *qdrant.ScrollPoints,
) ([]*qdrant.RetrievedPoint, *qdrant.PointId, error) {
	return c.points.ScrollAndOffset(ctx, request)
}

// ScrollAll returns an iterator over all matching points.
func (c *Client) ScrollAll(ctx context.Context, request *qdrant.ScrollPoints) *qdrant.ScrollIterator {
	return c.points.ScrollAll(ctx, request)
}

// Counts points.
func (c *Client) Count(ctx context.Context, request *qdrant.CountPoints) (uint64, error) {
	return c.points.Count(ctx, request)
}

// Inserts or overwrites points.
func (c *Client) Upsert(ctx context.Context, request *qdrant.UpsertPoints) (*qdrant.UpdateResult, error) {
	return c.points.Upsert(ctx, request)
}

// Updates vectors of points.
func (c *Client) UpdateVectors(ctx context.Context, request *qdrant.UpdatePointVectors) (*qdrant.UpdateResult, error) {
	return c.points.UpdateVectors(ctx, request)
}

// Removes vectors from points.
func (c *Client) DeleteVectors(ctx context.Context, request *qdrant.DeletePointVectors) (*qdrant.UpdateResult, error) {
	return c.points.DeleteVectors(ctx, request)
}

// Deletes points.
func (c *Client) Delete(ctx context.Context, request *qdrant.DeletePoints) (*qdrant.UpdateResult, error) {
	return c.points.Delete(ctx, request)
}

// Sets payload fields of points.
func (c *Client) SetPayload(ctx context.Context, request *qdrant.SetPayloadPoints) (*qdrant.UpdateResult, error) {
	return c.points.SetPayload(ctx, request)
}

// Replaces the entire payload of points.
func (c *Client) OverwritePayload(ctx context.Context, request *qdrant.SetPayloadPoints) (*qdrant.UpdateResult, error) {
	return c.points.OverwritePayload(ctx, request)
}

// Removes payload fields from points.
func (c *Client) DeletePayload(ctx context.Context, request *qdrant.DeletePayloadPoints) (*qdrant.UpdateResult, error) {
	return c.points.DeletePayload(ctx, request)
}

// Removes the entire payload of points.
func (c *Client) ClearPayload(ctx context.Context, request *qdrant.ClearPayloadPoints) (*qdrant.UpdateResult, error) {
	return c.points.ClearPayload(ctx, request)
}

// Runs several update operations in one request.
func (c *Client) UpdateBatch(ctx context.Context, request *qdrant.UpdateBatchPoints) ([]*qdrant.UpdateResult, error) {
	return c.points.UpdateBatch(ctx, request)
}
