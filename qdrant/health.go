package qdrant

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/connectivity"
)

// ConnectionState represents the health state of a connection
type ConnectionState int32

const (
	// ConnectionStateHealthy indicates the connection is working properly
	ConnectionStateHealthy ConnectionState = iota
	// ConnectionStateUnhealthy indicates the connection has failed health checks
	ConnectionStateUnhealthy
	// ConnectionStateRecovering indicates the connection is attempting to recover
	ConnectionStateRecovering
	// ConnectionStateDisconnected indicates the connection is permanently disconnected
	ConnectionStateDisconnected
)

func (s ConnectionState) String() string {
	switch s {
	case ConnectionStateHealthy:
		return "healthy"
	case ConnectionStateUnhealthy:
		return "unhealthy"
	case ConnectionStateRecovering:
		return "recovering"
	case ConnectionStateDisconnected:
		return "disconnected"
	default:
		return "unknown"
	}
}

// HealthCheckConfig configures health monitoring behavior
type HealthCheckConfig struct {
	// Interval between health checks
	Interval time.Duration
	// Timeout for individual health check requests
	Timeout time.Duration
	// Number of consecutive failures before marking connection as unhealthy
	FailureThreshold int
	// Number of consecutive successes before marking connection as healthy again
	RecoveryThreshold int
	// Whether to enable automatic recovery attempts
	EnableAutoRecovery bool
	// Interval between recovery attempts
	RecoveryInterval time.Duration
	// Maximum number of recovery attempts before giving up
	MaxRecoveryAttempts int
	// Logger for health check events
	Logger *slog.Logger
}

// DefaultHealthCheckConfig returns sensible defaults for health checking
func DefaultHealthCheckConfig() *HealthCheckConfig {
	return &HealthCheckConfig{
		Interval:            30 * time.Second,
		Timeout:             5 * time.Second,
		FailureThreshold:    3,
		RecoveryThreshold:   2,
		EnableAutoRecovery:  true,
		RecoveryInterval:    10 * time.Second,
		MaxRecoveryAttempts: 5,
		Logger:              slog.Default(),
	}
}

// ConnectionHealth tracks the health status of a single connection
type ConnectionHealth struct {
	client               *GrpcClient
	state                int32 // atomic ConnectionState
	consecutiveFailures  int32
	consecutiveSuccesses int32
	lastHealthCheck      time.Time
	lastError            error
	recoveryAttempts     int32
	mu                   sync.RWMutex
}

// NewConnectionHealth creates a new health tracker for a connection
func NewConnectionHealth(client *GrpcClient) *ConnectionHealth {
	return &ConnectionHealth{
		client:          client,
		state:           int32(ConnectionStateHealthy),
		lastHealthCheck: time.Now(),
	}
}

// GetState returns the current connection state
func (ch *ConnectionHealth) GetState() ConnectionState {
	return ConnectionState(atomic.LoadInt32(&ch.state))
}

// setState atomically updates the connection state
func (ch *ConnectionHealth) setState(state ConnectionState) {
	atomic.StoreInt32(&ch.state, int32(state))
}

// GetLastError returns the last error encountered during health checks
func (ch *ConnectionHealth) GetLastError() error {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.lastError
}

// setLastError sets the last error encountered
func (ch *ConnectionHealth) setLastError(err error) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.lastError = err
}

// GetLastHealthCheck returns the timestamp of the last health check
func (ch *ConnectionHealth) GetLastHealthCheck() time.Time {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return ch.lastHealthCheck
}

// setLastHealthCheck updates the last health check timestamp
func (ch *ConnectionHealth) setLastHealthCheck(t time.Time) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.lastHealthCheck = t
}

// IsHealthy returns true if the connection is currently healthy
func (ch *ConnectionHealth) IsHealthy() bool {
	state := ch.GetState()
	return state == ConnectionStateHealthy || state == ConnectionStateRecovering
}

// HealthMonitor manages health checking for all connections in the pool
type HealthMonitor struct {
	config      *HealthCheckConfig
	connections []*ConnectionHealth
	stopCh      chan struct{}
	stoppedCh   chan struct{}
	mu          sync.RWMutex
	running     int32
}

// NewHealthMonitor creates a new health monitor for the given connections
func NewHealthMonitor(clients []*GrpcClient, config *HealthCheckConfig) *HealthMonitor {
	if config == nil {
		config = DefaultHealthCheckConfig()
	}

	connections := make([]*ConnectionHealth, len(clients))
	for i, client := range clients {
		connections[i] = NewConnectionHealth(client)
	}

	return &HealthMonitor{
		config:      config,
		connections: connections,
		stopCh:      make(chan struct{}),
		stoppedCh:   make(chan struct{}),
	}
}

// Start begins health monitoring in the background
func (hm *HealthMonitor) Start() {
	if !atomic.CompareAndSwapInt32(&hm.running, 0, 1) {
		return // Already running
	}

	go hm.monitorLoop()
}

// Stop stops health monitoring
func (hm *HealthMonitor) Stop() {
	if !atomic.CompareAndSwapInt32(&hm.running, 1, 0) {
		return // Not running
	}

	close(hm.stopCh)
	<-hm.stoppedCh
}

// IsRunning returns true if the health monitor is currently running
func (hm *HealthMonitor) IsRunning() bool {
	return atomic.LoadInt32(&hm.running) == 1
}

// GetHealthyConnections returns a slice of healthy connection indices
func (hm *HealthMonitor) GetHealthyConnections() []int {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	var healthy []int
	for i, conn := range hm.connections {
		if conn.IsHealthy() {
			healthy = append(healthy, i)
		}
	}
	return healthy
}

// GetConnectionHealth returns the health status of a specific connection
func (hm *HealthMonitor) GetConnectionHealth(index int) *ConnectionHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	if index < 0 || index >= len(hm.connections) {
		return nil
	}
	return hm.connections[index]
}

// GetOverallHealth returns a summary of the pool's health
func (hm *HealthMonitor) GetOverallHealth() PoolHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	var healthy, unhealthy, recovering, disconnected int
	for _, conn := range hm.connections {
		switch conn.GetState() {
		case ConnectionStateHealthy:
			healthy++
		case ConnectionStateUnhealthy:
			unhealthy++
		case ConnectionStateRecovering:
			recovering++
		case ConnectionStateDisconnected:
			disconnected++
		}
	}

	return PoolHealth{
		TotalConnections:        len(hm.connections),
		HealthyConnections:      healthy,
		UnhealthyConnections:    unhealthy,
		RecoveringConnections:   recovering,
		DisconnectedConnections: disconnected,
	}
}

// PoolHealth represents the overall health status of the connection pool
type PoolHealth struct {
	TotalConnections        int
	HealthyConnections      int
	UnhealthyConnections    int
	RecoveringConnections   int
	DisconnectedConnections int
}

// IsHealthy returns true if at least one connection is healthy
func (ph PoolHealth) IsHealthy() bool {
	return ph.HealthyConnections > 0 || ph.RecoveringConnections > 0
}

// HealthRatio returns the ratio of healthy connections to total connections
func (ph PoolHealth) HealthRatio() float64 {
	if ph.TotalConnections == 0 {
		return 0
	}
	return float64(ph.HealthyConnections) / float64(ph.TotalConnections)
}

// monitorLoop runs the main health monitoring loop
func (hm *HealthMonitor) monitorLoop() {
	defer close(hm.stoppedCh)

	ticker := time.NewTicker(hm.config.Interval)
	defer ticker.Stop()

	// Perform initial health check
	hm.performHealthChecks()

	for {
		select {
		case <-hm.stopCh:
			return
		case <-ticker.C:
			hm.performHealthChecks()
		}
	}
}

// performHealthChecks checks the health of all connections
func (hm *HealthMonitor) performHealthChecks() {
	hm.mu.RLock()
	connections := make([]*ConnectionHealth, len(hm.connections))
	copy(connections, hm.connections)
	hm.mu.RUnlock()

	var wg sync.WaitGroup
	for _, conn := range connections {
		wg.Add(1)
		go func(ch *ConnectionHealth) {
			defer wg.Done()
			hm.checkConnectionHealth(ch)
		}(conn)
	}
	wg.Wait()
}

// checkConnectionHealth performs a health check on a single connection
func (hm *HealthMonitor) checkConnectionHealth(ch *ConnectionHealth) {
	ctx, cancel := context.WithTimeout(context.Background(), hm.config.Timeout)
	defer cancel()

	ch.setLastHealthCheck(time.Now())

	// Check gRPC connection state first
	grpcState := ch.client.Conn().GetState()
	if grpcState == connectivity.Shutdown || grpcState == connectivity.TransientFailure {
		hm.handleHealthCheckFailure(ch, fmt.Errorf("gRPC connection in state: %v", grpcState))
		return
	}

	// Perform actual health check using Qdrant's health endpoint
	_, err := ch.client.Qdrant().HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		hm.handleHealthCheckFailure(ch, err)
		return
	}

	hm.handleHealthCheckSuccess(ch)
}

// handleHealthCheckFailure processes a failed health check
func (hm *HealthMonitor) handleHealthCheckFailure(ch *ConnectionHealth, err error) {
	ch.setLastError(err)
	failures := atomic.AddInt32(&ch.consecutiveFailures, 1)
	atomic.StoreInt32(&ch.consecutiveSuccesses, 0)

	currentState := ch.GetState()

	hm.config.Logger.Warn("Health check failed",
		"error", err,
		"consecutive_failures", failures,
		"current_state", currentState.String())

	if failures >= int32(hm.config.FailureThreshold) && currentState == ConnectionStateHealthy {
		ch.setState(ConnectionStateUnhealthy)
		hm.config.Logger.Error("Connection marked as unhealthy",
			"consecutive_failures", failures,
			"threshold", hm.config.FailureThreshold)

		if hm.config.EnableAutoRecovery {
			go hm.attemptRecovery(ch)
		}
	}
}

// handleHealthCheckSuccess processes a successful health check
func (hm *HealthMonitor) handleHealthCheckSuccess(ch *ConnectionHealth) {
	ch.setLastError(nil)
	successes := atomic.AddInt32(&ch.consecutiveSuccesses, 1)
	atomic.StoreInt32(&ch.consecutiveFailures, 0)

	currentState := ch.GetState()

	if currentState != ConnectionStateHealthy && successes >= int32(hm.config.RecoveryThreshold) {
		ch.setState(ConnectionStateHealthy)
		atomic.StoreInt32(&ch.recoveryAttempts, 0)
		hm.config.Logger.Info("Connection recovered",
			"consecutive_successes", successes,
			"recovery_threshold", hm.config.RecoveryThreshold)
	}
}

// attemptRecovery tries to recover an unhealthy connection
func (hm *HealthMonitor) attemptRecovery(ch *ConnectionHealth) {
	if ch.GetState() != ConnectionStateUnhealthy {
		return
	}

	ch.setState(ConnectionStateRecovering)
	attempts := atomic.AddInt32(&ch.recoveryAttempts, 1)

	hm.config.Logger.Info("Starting connection recovery",
		"attempt", attempts,
		"max_attempts", hm.config.MaxRecoveryAttempts)

	ticker := time.NewTicker(hm.config.RecoveryInterval)
	defer ticker.Stop()

	for attempts <= int32(hm.config.MaxRecoveryAttempts) {
		select {
		case <-hm.stopCh:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), hm.config.Timeout)
			_, err := ch.client.Qdrant().HealthCheck(ctx, &HealthCheckRequest{})
			cancel()

			if err == nil {
				hm.config.Logger.Info("Connection recovery successful", "attempt", attempts)
				return // Success will be handled by the regular health check
			}

			attempts = atomic.AddInt32(&ch.recoveryAttempts, 1)
			hm.config.Logger.Warn("Connection recovery attempt failed",
				"attempt", attempts,
				"error", err)
		}
	}

	// Max attempts reached
	ch.setState(ConnectionStateDisconnected)
	hm.config.Logger.Error("Connection recovery failed after max attempts",
		"max_attempts", hm.config.MaxRecoveryAttempts)
}
