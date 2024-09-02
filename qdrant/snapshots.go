package qdrant

import (
	"context"
)

// Creates a snapshot of a specific collection.
// Snapshots are read-only copies of the collection data, which can be used for backup and restore purposes.
// The snapshot is created asynchronously and does not block the collection usage.
//
// Parameters:
//   - ctx: The context for the request
//   - collection: The name of the collection to create a snapshot for
//
// Returns:
//   - *SnapshotDescription: Description of the created snapshot
//   - error: Any error encountered during the snapshot creation
func (c *Client) CreateSnapshot(ctx context.Context, collection string) (*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().Create(ctx, &CreateSnapshotRequest{
		CollectionName: collection,
	})
	if err != nil {
		return nil, newQdrantErr(err, "CreateSnapshot")
	}
	return resp.GetSnapshotDescription(), nil
}

// Retrieves a list of all snapshots for a specific collection.
//
// Parameters:
//   - ctx: The context for the request
//   - collection: The name of the collection to list snapshots for
//
// Returns:
//   - []*SnapshotDescription: A slice of snapshot descriptions
//   - error: Any error encountered while listing snapshots
func (c *Client) ListSnapshots(ctx context.Context, collection string) ([]*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().List(ctx, &ListSnapshotsRequest{
		CollectionName: collection,
	})
	if err != nil {
		return nil, newQdrantErr(err, "ListSnapshots")
	}
	return resp.GetSnapshotDescriptions(), nil
}

// Removes a specific snapshot of a collection.
//
// Parameters:
//   - ctx: The context for the request
//   - collection: The name of the collection the snapshot belongs to
//   - snapshot: The name of the snapshot to delete
//
// Returns:
//   - error: Any error encountered while deleting the snapshot
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

// Creates a snapshot of the entire storage, including all collections.
// This operation is useful for creating full backups of the Qdrant instance.
//
// Parameters:
//   - ctx: The context for the request
//
// Returns:
//   - *SnapshotDescription: Description of the created full snapshot
//   - error: Any error encountered during the full snapshot creation
func (c *Client) CreateFullSnapshot(ctx context.Context) (*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().CreateFull(ctx, &CreateFullSnapshotRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "CreateFullSnapshot")
	}
	return resp.GetSnapshotDescription(), nil
}

// ListFullSnapshots retrieves a list of all full snapshots of the storage.
//
// Parameters:
//   - ctx: The context for the request
//
// Returns:
//   - []*SnapshotDescription: A slice of full snapshot descriptions
//   - error: Any error encountered while listing full snapshots
func (c *Client) ListFullSnapshots(ctx context.Context) ([]*SnapshotDescription, error) {
	resp, err := c.GetSnapshotsClient().ListFull(ctx, &ListFullSnapshotsRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "ListFullSnapshots")
	}
	return resp.GetSnapshotDescriptions(), nil
}

// Removes a specific full snapshot of the storage.
//
// Parameters:
//   - ctx: The context for the request
//   - snapshot: The name of the full snapshot to delete
//
// Returns:
//   - error: Any error encountered while deleting the full snapshot
func (c *Client) DeleteFullSnapshot(ctx context.Context, snapshot string) error {
	_, err := c.GetSnapshotsClient().DeleteFull(ctx, &DeleteFullSnapshotRequest{
		SnapshotName: snapshot,
	})
	if err != nil {
		return newQdrantErr(err, "DeleteFullSnapshot")
	}
	return nil
}
