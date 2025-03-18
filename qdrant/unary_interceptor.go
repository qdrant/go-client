package qdrant

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.ResourceExhausted {
			return err // return as is
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("retry-after"); len(values) > 0 {
				parsed, err := strconv.Atoi(values[0])
				if err == nil {
					return &QdrantResourceExhaustedError{
						err,
						parsed,
					}
				}
			}
		}

		return err
	}
}
