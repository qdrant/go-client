package qdrant_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qdrant/go-client/qdrant"
)

func TestClientHostAddress(t *testing.T) {
	for _, tt := range []struct {
		host string
		port int
		want string
	}{
		{host: "", want: "localhost:6334"},
		{host: "qdrant.example.com", port: 6335, want: "qdrant.example.com:6335"},
		{host: "127.0.0.1", want: "127.0.0.1:6334"},
		{host: "::1", want: "[::1]:6334"},
		{host: "[::1]", want: "[::1]:6334"},
		// The zone separator is percent-encoded in the gRPC target.
		{host: "fe80::1%eth0", want: "[fe80::1%25eth0]:6334"},
	} {
		client, err := qdrant.NewClient(&qdrant.Config{
			Host: tt.host, Port: tt.port, PoolSize: 1, SkipCompatibilityCheck: true,
		})
		require.NoError(t, err, tt.host)
		require.Equal(t, tt.want, client.GetGrpcClient().Conn().Target())
		require.NoError(t, client.Close())
	}
}
