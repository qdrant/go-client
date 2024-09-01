package qdrant

import (
	"context"
)

func (c *Client) CreateSnapshot(ctx context.Context, collection string) (*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().Create(ctx, &CreateSnapshotRequest{
		CollectionName: collection,
	})
	if err != nil {
		return nil, newQdrantErr(err, "CreateSnapshot")
	}
	return resp.GetSnapshotDescription(), nil
}

func (c *Client) ListSnapshots(ctx context.Context, collection string) ([]*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().List(ctx, &ListSnapshotsRequest{
		CollectionName: collection,
	})
	if err != nil {
		return nil, newQdrantErr(err, "ListSnapshots")
	}
	return resp.GetSnapshotDescriptions(), nil
}

func (c *Client) DeleteSnapshot(ctx context.Context, collection string, snapshot string) error {
	_, err := c.GetSnapshotsClient().Delete(ctx, &DeleteSnapshotRequest{
		CollectionName: collection,
		SnapshotName:   snapshot,
	})
	if err != nil {
		return newQdrantErr(err, "DeleteSnapshot")
	}
	return nil
}

func (c *Client) CreateFullSnapshot(ctx context.Context) (*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().CreateFull(ctx, &CreateFullSnapshotRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "CreateFullSnapshot")
	}
	return resp.GetSnapshotDescription(), nil
}

func (c *Client) ListFullSnapshots(ctx context.Context) ([]*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().ListFull(ctx, &ListFullSnapshotsRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "ListFullSnapshots")
	}
	return resp.GetSnapshotDescriptions(), nil
}

func (c *Client) DeleteFullSnapshot(ctx context.Context, snapshot string) error {
	_, err := c.GetSnapshotsClient().DeleteFull(ctx, &DeleteFullSnapshotRequest{
		SnapshotName: snapshot,
	})
	if err != nil {
		return newQdrantErr(err, "DeleteFullSnapshot")
	}
	return nil
}
