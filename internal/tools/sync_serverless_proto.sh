#!/usr/bin/env bash

set -euo pipefail

PROJECT_ROOT="$(pwd)/$(dirname "$0")/../../"
cd "$PROJECT_ROOT"

PROTO_FILE="internal/proto/serverless_collections.proto"

# buf.validate is server-side only and would need an extra proto dependency. Only plain
# single-option lists are stripped. Anything else is left in place and fails the check.
PROTO=$(curl -fsSL https://raw.githubusercontent.com/qdrant/qdrant-cloud-public-api/main/proto/qdrant/serverless/collections.proto |
    perl -0pe 's|import "buf/validate/validate.proto";\n||; s/\s*\[\(buf\.validate\.[^\[\]"\x27,]*\]//g')
if grep -qE 'buf[./]validate' <<< "$PROTO"; then
    echo "unexpected serverless proto, update $0" >&2
    exit 1
fi

cat > "$PROTO_FILE" <<EOF
// Source: https://github.com/qdrant/qdrant-cloud-public-api/blob/main/proto/qdrant/serverless/collections.proto
// Renamed because the protobuf registry already has the regular qdrant "collections.proto".
// buf.validate options are stripped: they are server-side only.
// Regenerate with: sh internal/tools/sync_serverless_proto.sh

$PROTO
EOF

./internal/tools/generate_proto_go.sh
