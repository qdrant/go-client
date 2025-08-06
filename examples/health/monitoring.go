package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	if err := runHealthMonitoringExample(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runHealthMonitoringExample() error {
	fmt.Println("Qdrant Health Monitoring Example")
	fmt.Println("=================================")

	// Create a client with health monitoring enabled
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:     "localhost",
		Port:     6334,
		PoolSize: 3, // Use multiple connections to demonstrate pool health
		HealthCheck: &qdrant.HealthCheckConfig{
			Interval:            2 * time.Second, // Check health every 2 seconds
			Timeout:             1 * time.Second, // 1 second timeout for health checks
			FailureThreshold:    2,               // Mark unhealthy after 2 consecutive failures
			RecoveryThreshold:   1,               // Mark healthy after 1 successful check
			EnableAutoRecovery:  true,            // Enable automatic recovery attempts
			RecoveryInterval:    1 * time.Second, // Wait 1 second between recovery attempts
			MaxRecoveryAttempts: 3,               // Try up to 3 times to recover
			Logger:              slog.Default(),  // Use default logger
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Check if health monitoring is enabled
	monitor := client.GetHealthMonitor()
	if monitor == nil {
		fmt.Println("Health monitoring is disabled")
		return fmt.Errorf("health monitoring not enabled")
	}

	fmt.Println("Health monitoring is enabled")
	fmt.Printf("Health monitor is running: %v\n", monitor.IsRunning())

	// Wait for initial health checks to complete
	fmt.Println("Waiting for initial health checks...")
	time.Sleep(1 * time.Second)

	// Check overall pool health
	poolHealth := client.GetPoolHealth()
	if poolHealth == nil {
		return fmt.Errorf("pool health not available")
	}

	fmt.Printf("Pool Health Status:\n")
	fmt.Printf("   Total connections: %d\n", poolHealth.TotalConnections)
	fmt.Printf("   Healthy connections: %d\n", poolHealth.HealthyConnections)
	fmt.Printf("   Unhealthy connections: %d\n", poolHealth.UnhealthyConnections)
	fmt.Printf("   Recovering connections: %d\n", poolHealth.RecoveringConnections)
	fmt.Printf("   Disconnected connections: %d\n", poolHealth.DisconnectedConnections)
	fmt.Printf("   Health ratio: %.2f%%\n", poolHealth.HealthRatio()*100)
	fmt.Printf("   Overall healthy: %v\n", poolHealth.IsHealthy())

	// Check individual connection health
	fmt.Println("\nIndividual Connection Health:")
	for i := 0; i < poolHealth.TotalConnections; i++ {
		connHealth := monitor.GetConnectionHealth(i)
		if connHealth != nil {
			fmt.Printf("   Connection %d: %s", i, connHealth.GetState().String())
			if lastErr := connHealth.GetLastError(); lastErr != nil {
				fmt.Printf(" (last error: %v)", lastErr)
			}
			if lastCheck := connHealth.GetLastHealthCheck(); !lastCheck.IsZero() {
				fmt.Printf(" (last check: %s)", lastCheck.Format("15:04:05"))
			}
			fmt.Println()
		}
	}

	// Test client operations
	fmt.Println("\nTesting client operations...")

	// Perform health check
	healthResult, err := client.HealthCheck(ctx)
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
	} else {
		fmt.Printf("Health check successful: %s\n", healthResult.GetVersion())
	}

	// Create a test collection
	collectionName := "health_test_collection"
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		fmt.Printf("Failed to create collection: %v\n", err)
	} else {
		fmt.Printf("Created collection: %s\n", collectionName)
	}

	// Insert some test data
	points := []*qdrant.PointStruct{
		{
			Id:      qdrant.NewIDNum(1),
			Vectors: qdrant.NewVectors(0.1, 0.2, 0.3, 0.4),
			Payload: qdrant.NewValueMap(map[string]any{"test": "data"}),
		},
		{
			Id:      qdrant.NewIDNum(2),
			Vectors: qdrant.NewVectors(0.2, 0.3, 0.4, 0.5),
			Payload: qdrant.NewValueMap(map[string]any{"test": "data2"}),
		},
	}

	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	if err != nil {
		fmt.Printf("Failed to upsert points: %v\n", err)
	} else {
		fmt.Printf("Upserted %d points\n", len(points))
	}

	// Perform a search
	searchResults, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(0.1, 0.2, 0.3, 0.4),
		Limit:          qdrant.PtrOf(uint64(5)),
	})
	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
	} else {
		fmt.Printf("Search successful, found %d results\n", len(searchResults))
	}

	// Demonstrate waiting for healthy connections
	fmt.Println("\nTesting WaitForHealthy...")
	waitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.WaitForHealthy(waitCtx)
	if err != nil {
		fmt.Printf("WaitForHealthy failed: %v\n", err)
	} else {
		fmt.Println("All connections are healthy!")
	}

	// Monitor health for a short period
	fmt.Println("\n📈 Monitoring health for 10 seconds...")
	monitorCtx, monitorCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer monitorCancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	monitorCount := 0
	for {
		select {
		case <-monitorCtx.Done():
			fmt.Println("Health monitoring demonstration complete")
			goto cleanup
		case <-ticker.C:
			monitorCount++
			currentHealth := client.GetPoolHealth()
			if currentHealth != nil {
				fmt.Printf("[%d] Current health: %d/%d healthy (%.1f%%) - %v\n",
					monitorCount,
					currentHealth.HealthyConnections,
					currentHealth.TotalConnections,
					currentHealth.HealthRatio()*100,
					currentHealth.IsHealthy())
			}
		}
	}

cleanup:
	// Clean up test collection
	err = client.DeleteCollection(ctx, collectionName)
	if err != nil {
		fmt.Printf("Failed to delete test collection: %v\n", err)
	} else {
		fmt.Printf("Cleaned up test collection: %s\n", collectionName)
	}

	// Final health status
	fmt.Println("\nFinal Health Status:")
	finalHealth := client.GetPoolHealth()
	if finalHealth != nil {
		fmt.Printf("   Total connections: %d\n", finalHealth.TotalConnections)
		fmt.Printf("   Healthy connections: %d\n", finalHealth.HealthyConnections)
		fmt.Printf("   Health ratio: %.2f%%\n", finalHealth.HealthRatio()*100)
		fmt.Printf("   Overall healthy: %v\n", finalHealth.IsHealthy())
	}

	fmt.Println("\nHealth monitoring example completed successfully!")
	return nil
}
