package qdrant

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"net/netip"
	"strconv"
	"strings"
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
	apiKeyHeader               = "api-key"
	defaultHost                = "localhost"
	defaultPort                = 6334
	defaultVersionCheckTimeout = time.Minute
	defaultMaxRecvMsgSize      = 64 << 20 // 64MB
)

var errNilConfig = errors.New("qdrant: config must not be nil")

// Configuration options for the client.
type Config struct {
	// Hostname or IP address of the Qdrant server. Defaults to "localhost".
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
	// RetryConfig enables automatic retries with exponential backoff for
	// transient gRPC errors (ResourceExhausted, Unavailable).
	// If nil, no automatic retries are performed.
	RetryConfig *RetryConfig
	// VersionCheckTimeout specifies the timeout used when probing the server for its
	// version during client construction (compatibility check).
	// If 0, defaults to 1 minute.
	VersionCheckTimeout time.Duration
	// Headers specifies optional headers to send with every gRPC request.
	Headers map[string]string
	// Logger used for client-side warnings such as a failed server version check.
	// If nil, slog.Default() is used.
	Logger *slog.Logger
	// MaxRecvMsgSize is the largest response, in bytes, the client accepts.
	// Larger responses fail with a ResourceExhausted error.
	// If 0, defaults to 64MB.
	// A grpc.MaxCallRecvMsgSize call option in GrpcOptions takes precedence.
	MaxRecvMsgSize uint
}

// Internal method.
func (c *Config) getMaxRecvMsgSize() int {
	if c.MaxRecvMsgSize == 0 {
		return defaultMaxRecvMsgSize
	}
	// gRPC takes an int, so cap the value to avoid wrapping to a negative limit.
	return int(min(c.MaxRecvMsgSize, uint(math.MaxInt)))
}

// Internal method.
func (c *Config) getLogger() *slog.Logger {
	if c.Logger != nil {
		return c.Logger
	}
	return slog.Default()
}

// Internal method.
func (c *Config) getVersionCheckTimeout() time.Duration {
	if c.VersionCheckTimeout > 0 {
		return c.VersionCheckTimeout
	}
	return defaultVersionCheckTimeout
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
	// netip.ParseAddr accepts bare IPs, including IPv6 addresses with a zone
	// identifier (for example "fe80::1%eth0"). gRPC parses the target as a URI,
	// so the zone separator must be percent-encoded.
	if _, err := netip.ParseAddr(host); err == nil {
		return strings.ReplaceAll(net.JoinHostPort(host, strconv.Itoa(port)), "%", "%25")
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
		c.getLogger().Warn("API key is being used without TLS(HTTPS). It will be transmitted in plaintext.")
	}
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

// Internal method.
//
//nolint:lll
func (c *Config) getMetadataInterceptor() grpc.DialOption {
	var kvPairs []string
	if c.APIKey != "" {
		kvPairs = append(kvPairs, apiKeyHeader, c.APIKey)
	}
	for k, v := range c.Headers {
		kvPairs = append(kvPairs, k, v)
	}
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := metadata.AppendToOutgoingContext(ctx, kvPairs...)
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
			parsed, parseErr := strconv.Atoi(values[0])
			if parseErr == nil {
				return &QdrantResourceExhaustedError{
					Reason:      st.Message(),
					RetryAfterS: parsed,
					err:         err,
				}
			}
			return errors.Join(fmt.Errorf("parse retry-after header %q: %w", values[0], parseErr), err)
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
