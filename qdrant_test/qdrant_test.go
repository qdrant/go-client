package qdrant_test

import (
	"context"
	"testing"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	apiKey := "<HEALTHCHECK_TEST>"

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

	resp, err := client.HealthCheck(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
