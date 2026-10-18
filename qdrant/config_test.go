package qdrant_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/qdrant/go-client/qdrant"
)

type healthServer struct {
	qdrant.UnimplementedQdrantServer
}

func (healthServer) HealthCheck(context.Context, *qdrant.HealthCheckRequest) (*qdrant.HealthCheckReply, error) {
	return &qdrant.HealthCheckReply{Title: "host-address-test"}, nil
}

func TestClientHostAddress(t *testing.T) {
	for _, tt := range []struct {
		name    string
		network string
		address string
		host    string
	}{
		{name: "IPv4", network: "tcp4", address: "127.0.0.1:0", host: "127.0.0.1"},
		{name: "IPv6", network: "tcp6", address: "[::1]:0", host: "::1"},
		{name: "bracketed IPv6", network: "tcp6", address: "[::1]:0", host: "[::1]"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var lc net.ListenConfig
			listener, err := lc.Listen(context.Background(), tt.network, tt.address)
			require.NoError(t, err)
			server := grpc.NewServer()
			qdrant.RegisterQdrantServer(server, healthServer{})
			serveDone := make(chan error, 1)
			go func() { serveDone <- server.Serve(listener) }()
			t.Cleanup(func() {
				server.Stop()
				require.NoError(t, <-serveDone)
			})

			addr, ok := listener.Addr().(*net.TCPAddr)
			require.True(t, ok)
			client, err := qdrant.NewClient(&qdrant.Config{
				Host: tt.host, Port: addr.Port, PoolSize: 1, SkipCompatibilityCheck: true,
				GrpcOptions: []grpc.DialOption{grpc.WithNoProxy()},
			})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			response, err := client.HealthCheck(ctx)
			require.NoError(t, err)
			require.Equal(t, "host-address-test", response.GetTitle())
		})
	}
}
