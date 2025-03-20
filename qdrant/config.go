package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strconv"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
func (c *Config) getTransportCreds() grpc.DialOption {
	if c.UseTLS {
		if c.TLSConfig == nil {
			return grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS13,
			}))
		}
		return grpc.WithTransportCredentials(credentials.NewTLS(c.TLSConfig))
	} else if c.APIKey != "" {
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
