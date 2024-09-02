package qdrant_test

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestSnapshotsClient(t *testing.T) {
	collectionName := t.Name()
	apiKey := "<SNAPSHOTS_TEST>"

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	container, err := standaloneQdrant(ctx, apiKey)
	require.NoError(t, err)

	err = container.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	})
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6334/tcp")
	require.NoError(t, err)

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   port.Int(),
		APIKey: apiKey,
	})
	require.NoError(t, err)

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	require.NoError(t, err)

	t.Run("CreateSnapshot", func(t *testing.T) {
		snapshot, err := client.CreateSnapshot(ctx, collectionName)
		require.NoError(t, err)
		require.NotNil(t, snapshot)
	})

	t.Run("ListSnapshots", func(t *testing.T) {
		snapshots, err := client.ListSnapshots(ctx, collectionName)
		require.NoError(t, err)
		require.NotEmpty(t, snapshots)
	})

	t.Run("DeleteSnapshot", func(t *testing.T) {
		snapshots, err := client.ListSnapshots(ctx, collectionName)
		require.NoError(t, err)
		require.NotEmpty(t, snapshots)

		err = client.DeleteSnapshot(ctx, collectionName, snapshots[0].GetName())
		require.NoError(t, err)
	})

	t.Run("CreateFullSnapshot", func(t *testing.T) {
		snapshot, err := client.CreateFullSnapshot(ctx)
		require.NoError(t, err)
		require.NotNil(t, snapshot)
	})

	t.Run("ListFullSnapshots", func(t *testing.T) {
		snapshots, err := client.ListFullSnapshots(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, snapshots)
	})

	t.Run("DeleteFullSnapshot", func(t *testing.T) {
		snapshots, err := client.ListFullSnapshots(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, snapshots)

		err = client.DeleteFullSnapshot(ctx, snapshots[0].GetName())
		require.NoError(t, err)
	})

	err = client.DeleteCollection(ctx, collectionName)
	require.NoError(t, err)
}
