package qdrant_test

import (
	"testing"

	"github.com/qdrant/go-client/qdrant"
)

// Reusing one Config to construct multiple clients must not mutate the
// caller's GrpcOptions. NewGrpcClient used to write the combined dial
// options back into Config.GrpcOptions, so every subsequent client
// duplicated all default options (including interceptors).
func TestNewGrpcClientDoesNotMutateConfigGrpcOptions(t *testing.T) {
	config := &qdrant.Config{SkipCompatibilityCheck: true}

	_, err := qdrant.NewGrpcClient(config)
	if err != nil {
		t.Fatalf("first client: %v", err)
	}
	_, err = qdrant.NewGrpcClient(config)
	if err != nil {
		t.Fatalf("second client: %v", err)
	}

	if len(config.GrpcOptions) != 0 {
		t.Errorf("config.GrpcOptions mutated: got %d options, want 0", len(config.GrpcOptions))
	}
}
