package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

const (
	apiKeyHeader = "api-key"
	defaultHost  = "localhost"
	defaultPort  = 6334
)

// Configuration options for the client.
type Config struct {
	// Hostname of the Qdrant server. Defaults to "localhost".
	Host string
	// gRPC port of the Qdrant server. Defaults to 6334.
	Port int
	// API key to use for authentication. Defaults to "".
	APIKey string
	// Whether to use TLS for the connection. Defaults to false.
	UseTLS bool
	// TLS configuration to use for the connection.
	// If not provided, uses default config with minimum TLS version set to 1.3
	TLSConfig *tls.Config
	// Additional gRPC options to use for the connection.
	GrpcOptions []grpc.DialOption
	// Whether to check compatibility between server's version and client's. Defaults to false.
	SkipCompatibilityCheck bool
	// PoolSize specifies the number of connections to create.
	// If 0, the default of 3 will be used.
	// If 1 a single connection is used (aka no pool).
	// If greater than 1, a pool of connections is created and requests are distributed in a round-robin fashion.
	PoolSize uint
	// KeepAliveTime specifies the duration after which if the client does not see any activity (in seconds),
	// it pings the server to check if the transport is still alive.
	// If 0, the default is 10 seconds.
	// If set to -1, keepalive is disabled.
	KeepAliveTime int
	// KeepAliveTimeout specifies the duration the client waits for a response from the server after
	// sending a ping (in seconds).
	// If the server does not respond within this timeout, the connection is closed.
	// If set to 0, defaults to 2 seconds.
	// This setting is only used if keepalive is active (see KeepAliveTime).
	KeepAliveTimeout uint
}

// Internal method.
func (c *Config) getAddr() string {
	host := c.Host
	if host == "" {
		host = defaultHost
	}
	port := c.Port
	if port == 0 {
		port = defaultPort
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// Internal method.
func (c *Config) getKeepAliveParams() []grpc.DialOption {
	if c.KeepAliveTime == -1 {
		// Disabled
		return nil
	}
	// Default to 10 seconds
	keepAliveTime := 10
	if c.KeepAliveTime > 0 {
		keepAliveTime = c.KeepAliveTime
	}
	keepAliveTimeout := 2
	if c.KeepAliveTimeout > 0 {
		keepAliveTimeout = int(c.KeepAliveTimeout)
	}
	return []grpc.DialOption{getClientKeepAliveParams(keepAliveTime, keepAliveTimeout)}
}

// Internal method.
func (c *Config) getTransportCreds() grpc.DialOption {
	if c.UseTLS {
		if c.TLSConfig == nil {
			return grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS13,
			}))
		}
		return grpc.WithTransportCredentials(credentials.NewTLS(c.TLSConfig))
	} else if c.APIKey != "" {
		//nolint:noctx // We don't have context here.
		slog.Default().Warn("API key is being used without TLS(HTTPS). It will be transmitted in plaintext.")
	}
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

// Internal method.
//
//nolint:lll
func (c *Config) getAPIKeyInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := ctx
		if c.APIKey != "" {
			newCtx = metadata.AppendToOutgoingContext(ctx, apiKeyHeader, c.APIKey)
		}
		return invoker(newCtx, method, req, reply, cc, opts...)
	})
}

// Internal method.
func (c *Config) getRateLimitInterceptor() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		var md metadata.MD
		opts = append(opts, grpc.Trailer(&md))
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.ResourceExhausted {
			return err
		}
		if values := md.Get("retry-after"); len(values) > 0 {
			parsed, err := strconv.Atoi(values[0])
			if err == nil {
				return &QdrantResourceExhaustedError{
					st.Message(),
					parsed,
				}
			}
		}
		return err
	})
}

// Internal method.
func getClientKeepAliveParams(keepAliveTime, keepAliveTimeout int) grpc.DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// Send pings every keepAliveTime (default 10s) if no activity
		Time: time.Duration(keepAliveTime) * time.Second,
		// Wait keepAliveTimeout (default 2s) for ping ack before closing
		Timeout: time.Duration(keepAliveTimeout) * time.Second,
		// Send pings even with no active streams
		PermitWithoutStream: true,
	})
}
