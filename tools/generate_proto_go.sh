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

# SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROTO_DIR=./proto
OUT_DIR=./qdrant
PACKAGE_NAME="github.com/qdrant/go-client;qdrant"

protoc \
    --experimental_allow_proto3_optional \
    --proto_path=$PROTO_DIR/ \
    --go_out=$OUT_DIR \
    --go-grpc_out=$OUT_DIR \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    \
    --go_opt=Mcollections_service.proto=$PACKAGE_NAME \
    --go_opt=Mcollections.proto=$PACKAGE_NAME \
    --go_opt=Mjson_with_int.proto=$PACKAGE_NAME \
    --go_opt=Mpoints_service.proto=$PACKAGE_NAME \
    --go_opt=Mpoints.proto=$PACKAGE_NAME \
    --go_opt=Msnapshots_service.proto=$PACKAGE_NAME \
    --go_opt=Mqdrant.proto=$PACKAGE_NAME \
    \
    --go-grpc_opt=Mcollections_service.proto=$PACKAGE_NAME \
    --go-grpc_opt=Mcollections.proto=$PACKAGE_NAME \
    --go-grpc_opt=Mjson_with_int.proto=$PACKAGE_NAME \
    --go-grpc_opt=Mpoints_service.proto=$PACKAGE_NAME \
    --go-grpc_opt=Mpoints.proto=$PACKAGE_NAME \
    --go-grpc_opt=Msnapshots_service.proto=$PACKAGE_NAME \
    --go-grpc_opt=Mqdrant.proto=$PACKAGE_NAME \
    \
    $PROTO_DIR/collections_service.proto \
    $PROTO_DIR/collections.proto \
    $PROTO_DIR/json_with_int.proto \
    $PROTO_DIR/points_service.proto \
    $PROTO_DIR/points.proto \
    $PROTO_DIR/snapshots_service.proto \
    $PROTO_DIR/qdrant.proto \
