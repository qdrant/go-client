package qdrant_test

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestCollectionsClient(t *testing.T) {
	collectionName := t.Name()
	vectorSize := uint64(384)
	distance := qdrant.Distance_Cosine
	apiKey := "<COLLECTIONS_TEST>"

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	container, err := distributedQdrant(ctx, apiKey)
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

	t.Run("CreateCollection", func(t *testing.T) {
		err := client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     vectorSize,
				Distance: distance,
			}),
			ShardingMethod: qdrant.ShardingMethod_Custom.Enum(),
		})
		require.NoError(t, err)

		_, err = client.GetCollectionInfo(ctx, collectionName)
		require.NoError(t, err)
	})

	t.Run("CollectionExists", func(t *testing.T) {
		exists, err := client.CollectionExists(ctx, collectionName)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("GetCollection", func(t *testing.T) {
		collInfo, err := client.GetCollectionInfo(ctx, collectionName)
		require.NoError(t, err)
		require.Zero(t, collInfo.GetPointsCount())
	})

	t.Run("ListCollections", func(t *testing.T) {
		collections, err := client.ListCollections(ctx)
		require.NoError(t, err)
		require.Contains(t, collections, collectionName)
	})

	t.Run("UpdateCollection", func(t *testing.T) {
		threshold := uint64(1000)
		err := client.UpdateCollection(ctx, &qdrant.UpdateCollection{
			CollectionName: collectionName,
			OptimizersConfig: &qdrant.OptimizersConfigDiff{
				IndexingThreshold: &threshold,
			},
		})
		require.NoError(t, err)

		collInfo, err := client.GetCollectionInfo(ctx, collectionName)
		require.NoError(t, err)
		require.Equal(t, threshold, collInfo.GetConfig().GetOptimizerConfig().GetIndexingThreshold())
	})

	t.Run("AliasOperations", func(t *testing.T) {
		aliasName := "test_alias"
		newAliasName := "new_test_alias"

		t.Run("CreateAlias", func(t *testing.T) {
			err := client.CreateAlias(ctx, aliasName, collectionName)
			require.NoError(t, err)

			aliases, err := client.ListCollectionAliases(ctx, collectionName)
			require.NoError(t, err)
			require.Contains(t, aliases, aliasName)
		})

		t.Run("ListCollectionAliases", func(t *testing.T) {
			aliases, err := client.ListCollectionAliases(ctx, collectionName)
			require.NoError(t, err)
			require.Contains(t, aliases, aliasName)
		})

		t.Run("ListAliases", func(t *testing.T) {
			allAliases, err := client.ListAliases(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, allAliases)
		})

		t.Run("RenameAlias", func(t *testing.T) {
			err := client.RenameAlias(ctx, aliasName, newAliasName)
			require.NoError(t, err)

			aliases, err := client.ListCollectionAliases(ctx, collectionName)
			require.NoError(t, err)
			require.Contains(t, aliases, newAliasName)
			require.NotContains(t, aliases, aliasName)
		})

		t.Run("DeleteAlias", func(t *testing.T) {
			err := client.DeleteAlias(ctx, newAliasName)
			require.NoError(t, err)

			aliases, err := client.ListCollectionAliases(ctx, collectionName)
			require.NoError(t, err)
			require.NotContains(t, aliases, newAliasName)
		})
	})

	t.Run("ShardKeyOperations", func(t *testing.T) {
		shardKey := "test_shard_key"

		t.Run("CreateShardKey", func(t *testing.T) {
			err := client.CreateShardKey(ctx, collectionName, &qdrant.CreateShardKey{
				ShardKey: qdrant.NewShardKey(shardKey),
			})
			require.NoError(t, err)
		})

		t.Run("DeleteShardKey", func(t *testing.T) {
			err := client.DeleteShardKey(ctx, collectionName, &qdrant.DeleteShardKey{
				ShardKey: qdrant.NewShardKey(shardKey),
			})
			require.NoError(t, err)
		})
	})

	t.Run("DeleteCollection", func(t *testing.T) {
		err := client.DeleteCollection(ctx, collectionName)
		require.NoError(t, err)

		exists, err := client.CollectionExists(ctx, collectionName)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("CreateCollectionWithInvalidParams", func(t *testing.T) {
		err := client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: "",
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     0,
				Distance: qdrant.Distance_Cosine,
			}),
		})
		require.Error(t, err)
	})

	t.Run("UpdateNonExistentCollection", func(t *testing.T) {
		err := client.UpdateCollection(ctx, &qdrant.UpdateCollection{
			CollectionName: "non_existent_collection",
			OptimizersConfig: &qdrant.OptimizersConfigDiff{
				IndexingThreshold: new(uint64),
			},
		})
		require.Error(t, err)
	})

	t.Run("DeleteNonExistentCollection", func(t *testing.T) {
		err := client.DeleteCollection(ctx, "non_existent_collection")
		require.Error(t, err)
	})
}
