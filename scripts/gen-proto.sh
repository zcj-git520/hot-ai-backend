#!/bin/bash

# Protobuf 生成脚本

set -e

PROTO_DIR="pkg/proto"
PROTO_FILES=$(find . -name "*.proto" -type f)

echo "Generating protobuf code..."

for proto_file in $PROTO_FILES; do
    echo "Processing: $proto_file"
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           "$proto_file"
done

echo "Protobuf generation completed!"
