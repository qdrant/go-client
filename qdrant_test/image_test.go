package qdrant_test

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const TestImage string = "qdrant/qdrant:v1.16.1"

// We use an instance with distributed mode enabled
// to test methods like CreateShardKey(), DeleteShardKey().
func distributedQdrant(ctx context.Context, apiKey string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        TestImage,
		ExposedPorts: []string{"6334/tcp"},
		Env: map[string]string{
			"QDRANT__CLUSTER__ENABLED": "true",
			"QDRANT__SERVICE__API_KEY": apiKey,
		},
		Cmd: []string{"./qdrant", "--uri", "http://qdrant_node_1:6335"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6334/tcp").WithStartupTimeout(5 * time.Second),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	return container, err
}

func standaloneQdrant(ctx context.Context, apiKey string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        TestImage,
		ExposedPorts: []string{"6334/tcp"},
		Env: map[string]string{
			"QDRANT__SERVICE__API_KEY": apiKey,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6334/tcp").WithStartupTimeout(5 * time.Second),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})

	return container, err
}
