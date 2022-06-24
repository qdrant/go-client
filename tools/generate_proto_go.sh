#!/bin/bash

PROJECT_ROOT="$(pwd)/$(dirname "$0")/../"

QDRANT_PROTO_DIR='proto'

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# go install github.com/golang/protobuf/protoc-gen-go@v1.5.2


GOPATH=$(go env GOPATH)

case ":$PATH:" in
    *":$GOPATH/bin:"*) ;;
    *) export PATH="$GOPATH/bin:$PATH";;
esac


protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/qdrant.proto


protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/collections.proto

  protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/collections_service.proto

  protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/json_with_int.proto

  protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/points.proto

  protoc \
  --go_out=qdrant \
  --go_opt="Mqdrant.proto=./qdrant" \
  --go_opt="Mcollections.proto=./collections" \
  --go_opt="Mcollections_service.proto=./collections_service" \
  --go_opt="Mjson_with_int.proto=./json_with_int" \
  --go_opt="Mpoints.proto=./points" \
  --go_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_opt="Mqdrant.proto=./qdrant;qdrant" \
  --go-grpc_opt="Mcollections.proto=./collections" \
  --go-grpc_opt="Mcollections_service.proto=./collections_service" \
  --go-grpc_opt="Mjson_with_int.proto=./json_with_int" \
  --go-grpc_opt="Mpoints.proto=./points" \
  --go-grpc_opt="Mpoints_service.proto=./points_service" \
  --go-grpc_out=paths=source_relative:qdrant \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/points_service.proto
