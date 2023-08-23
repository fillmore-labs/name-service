#!/bin/sh -eu

if ! command -v protoc-gen-doc >/dev/null 2>&1; then
  go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
fi

protoc \
  --proto_path=proto \
  --go_out=api \
  --go_opt=paths=import \
  --go-grpc_out=api \
  --go-grpc_opt=paths=import \
  --doc_out=doc \
  --doc_opt=markdown,service.md \
  proto/fillmore_labs/*/*/*.proto
