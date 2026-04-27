package qdrant_test

import (
	"context"
	"testing"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestRetryConfig(t *testing.T) {
	apiKey := "<RETRY_TEST>"

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

	t.Run("WithRetryConfig", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
			RetryConfig: &qdrant.RetryConfig{
				MaxRetries:  3,
				BaseBackoff: 100 * time.Millisecond,
				MaxBackoff:  time.Second,
			},
		})
		require.NoError(t, err)

		resp, err := client.HealthCheck(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("WithRetryConfigDefaults", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
			RetryConfig: &qdrant.RetryConfig{
				MaxRetries: 3,
			},
		})
		require.NoError(t, err)

		resp, err := client.HealthCheck(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("WithoutRetryConfig", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
		})
		require.NoError(t, err)

		resp, err := client.HealthCheck(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("WithZeroMaxRetries", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
			RetryConfig: &qdrant.RetryConfig{
				MaxRetries: 0,
			},
		})
		require.NoError(t, err)

		resp, err := client.HealthCheck(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
