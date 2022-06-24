#!/bin/bash

PROJECT_ROOT="$(pwd)/$(dirname "$0")/../"

QDRANT_PROTO_DIR='proto'

# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# go install github.com/golang/protobuf/protoc-gen-go@v1.5.2


# protoc --proto_path=${QDRANT_PROTO_DIR} \
#     --go_opt="Mmilvus.proto=github.com/milvus-io/milvus-sdk-go/v2/internal/proto/server;server" \
#     --go_opt=Mcommon.proto=github.com/milvus-io/milvus-sdk-go/v2/internal/proto/common \
#     --go_opt=Mschema.proto=github.com/milvus-io/milvus-sdk-go/v2/internal/proto/schema \

#     --go_out=plugins=grpc,paths=source_relative:${PROTO_DIR}/server ${QDRANT_PROTO_DIR}/qdrant.proto

GOPATH=$(go env GOPATH)

case ":$PATH:" in
    *":$GOPATH/bin:"*) ;;
    *) export PATH="$GOPATH/bin:$PATH";;
esac


protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/qdrant.proto

protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/collections.proto

protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/points.proto

protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/collections_service.proto

protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/points_service.proto

protoc \
  --go_out=qdrant_client \
  --go_opt="Mqdrant.proto=/" \
  --go-grpc_out=paths=source_relative:qdrant_client \
  --proto_path=${QDRANT_PROTO_DIR} ${QDRANT_PROTO_DIR}/json_with_int.proto

