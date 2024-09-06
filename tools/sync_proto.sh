#!/usr/bin/env bash

set -e

BRANCH=${BRANCH:-"master"}
PROJECT_ROOT="$(pwd)/$(dirname "$0")/../"

cd $(mktemp -d)

git clone --sparse --branch $BRANCH --filter=blob:none --depth=1 https://github.com/qdrant/qdrant
cd qdrant
git sparse-checkout add lib/api/src/grpc/proto

PROTO_DIR="$(pwd)/lib/api/src/grpc/proto"

# Ensure current path is project root
cd $PROJECT_ROOT

CLIENT_DIR="proto"

cp $PROTO_DIR/*.proto $CLIENT_DIR/

# Remove internal service proto files and their imports
rm $CLIENT_DIR/collections_internal_service.proto
rm $CLIENT_DIR/points_internal_service.proto
rm $CLIENT_DIR/qdrant_internal_service.proto
rm $CLIENT_DIR/raft_service.proto
rm $CLIENT_DIR/health_check.proto
rm $CLIENT_DIR/shard_snapshots_service.proto
sed -i '
    /collections_internal_service.proto/d;
    /points_internal_service.proto/d;
    /qdrant_internal_service.proto/d;
    /raft_service.proto/d;
    /health_check.proto/d;
    /shard_snapshots_service.proto/d;
    ' $CLIENT_DIR/qdrant.proto

# Remove csharp option from proto files
sed -i '/option csharp_namespace = .*/d' $CLIENT_DIR/*.proto

sh tools/generate_proto_go.sh
