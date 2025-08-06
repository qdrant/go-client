package qdrant_test

import (
	"context"
	"testing"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/require"
)

func TestHealthMonitoring(t *testing.T) {
	// Enable health monitoring tests
	apiKey := "<HEALTH_TEST>"

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

	t.Run("BasicHealthCheck", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
		})
		require.NoError(t, err)
		defer client.Close()

		// Test basic health check functionality
		_, err = client.HealthCheck(ctx)
		require.NoError(t, err)
	})

	t.Run("ClientWithoutHealthMonitoring", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:   host,
			Port:   port.Int(),
			APIKey: apiKey,
			// No HealthCheck config
		})
		require.NoError(t, err)
		defer client.Close()

		// Verify health monitor is not initialized
		monitor := client.GetHealthMonitor()
		require.Nil(t, monitor)

		// Pool health should be nil
		poolHealth := client.GetPoolHealth()
		require.Nil(t, poolHealth)

		// IsHealthy should return true (assumes healthy)
		require.True(t, client.IsHealthy())

		// WaitForHealthy should return immediately
		err = client.WaitForHealthy(ctx)
		require.NoError(t, err)
	})

	t.Run("ClientWithHealthMonitoring", func(t *testing.T) {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host:     host,
			Port:     port.Int(),
			APIKey:   apiKey,
			PoolSize: 2,
			HealthCheck: &qdrant.HealthCheckConfig{
				Interval:            1 * time.Second,
				Timeout:             2 * time.Second,
				FailureThreshold:    2,
				RecoveryThreshold:   1,
				EnableAutoRecovery:  true,
				RecoveryInterval:    500 * time.Millisecond,
				MaxRecoveryAttempts: 3,
			},
		})
		require.NoError(t, err)
		defer client.Close()

		// Verify health monitor is initialized
		monitor := client.GetHealthMonitor()
		require.NotNil(t, monitor)
		require.True(t, monitor.IsRunning())

		// Wait for initial health checks
		time.Sleep(500 * time.Millisecond)

		// Pool health should be available
		poolHealth := client.GetPoolHealth()
		require.NotNil(t, poolHealth)
		require.Equal(t, 2, poolHealth.TotalConnections)
		require.True(t, poolHealth.IsHealthy())

		// Test individual connection health
		healthyConnections := monitor.GetHealthyConnections()
		require.NotEmpty(t, healthyConnections)

		// Test WaitForHealthy
		err = client.WaitForHealthy(ctx)
		require.NoError(t, err)
	})
}

func TestHealthCheckConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		config := qdrant.DefaultHealthCheckConfig()
		require.NotNil(t, config)
		require.Equal(t, 30*time.Second, config.Interval)
		require.Equal(t, 5*time.Second, config.Timeout)
		require.Equal(t, 3, config.FailureThreshold)
		require.Equal(t, 2, config.RecoveryThreshold)
		require.True(t, config.EnableAutoRecovery)
		require.Equal(t, 10*time.Second, config.RecoveryInterval)
		require.Equal(t, 5, config.MaxRecoveryAttempts)
		require.NotNil(t, config.Logger)
	})
}

func TestConnectionHealth(t *testing.T) {
	t.Run("ConnectionStates", func(t *testing.T) {
		// Test connection state string representations
		require.Equal(t, "healthy", qdrant.ConnectionStateHealthy.String())
		require.Equal(t, "unhealthy", qdrant.ConnectionStateUnhealthy.String())
		require.Equal(t, "recovering", qdrant.ConnectionStateRecovering.String())
		require.Equal(t, "disconnected", qdrant.ConnectionStateDisconnected.String())
	})

	t.Run("PoolHealthMethods", func(t *testing.T) {
		// Test PoolHealth methods
		poolHealth := qdrant.PoolHealth{
			TotalConnections:        4,
			HealthyConnections:      3,
			UnhealthyConnections:    1,
			RecoveringConnections:   0,
			DisconnectedConnections: 0,
		}

		require.True(t, poolHealth.IsHealthy())
		require.Equal(t, 0.75, poolHealth.HealthRatio())

		// Test unhealthy pool
		unhealthyPool := qdrant.PoolHealth{
			TotalConnections:        2,
			HealthyConnections:      0,
			UnhealthyConnections:    1,
			RecoveringConnections:   0,
			DisconnectedConnections: 1,
		}

		require.False(t, unhealthyPool.IsHealthy())
		require.Equal(t, 0.0, unhealthyPool.HealthRatio())

		// Test pool with recovering connections
		recoveringPool := qdrant.PoolHealth{
			TotalConnections:        2,
			HealthyConnections:      0,
			UnhealthyConnections:    0,
			RecoveringConnections:   2,
			DisconnectedConnections: 0,
		}

		require.True(t, recoveringPool.IsHealthy()) // Recovering counts as healthy
	})
}

func TestHealthMonitoringOperations(t *testing.T) {
	apiKey := "<HEALTH_TEST>"

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
		Host:     host,
		Port:     port.Int(),
		APIKey:   apiKey,
		PoolSize: 2,
		HealthCheck: &qdrant.HealthCheckConfig{
			Interval:            1 * time.Second,
			Timeout:             2 * time.Second,
			FailureThreshold:    2,
			RecoveryThreshold:   1,
			EnableAutoRecovery:  true,
			RecoveryInterval:    500 * time.Millisecond,
			MaxRecoveryAttempts: 3,
		},
	})
	require.NoError(t, err)
	defer client.Close()

	// Wait for health checks to initialize
	time.Sleep(500 * time.Millisecond)

	collectionName := "health_test_operations"

	// Create collection
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	require.NoError(t, err)

	// Cleanup
	defer func() {
		client.DeleteCollection(ctx, collectionName)
	}()

	// Test operations with health monitoring
	points := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
			Payload: qdrant.NewValueMap(map[string]any{"test": "data"}),
		},
	}

	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	require.NoError(t, err)

	// Verify health during operations
	poolHealth := client.GetPoolHealth()
	require.NotNil(t, poolHealth)
	require.True(t, poolHealth.IsHealthy())
}
