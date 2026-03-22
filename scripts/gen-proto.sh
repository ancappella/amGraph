#!/usr/bin/env bash
# Generate protobuf Go code without GNU make (works when macOS Command Line Tools / xcrun is broken).
# Root fix for "invalid active developer path" / missing xcrun:
#   sudo rm -rf /Library/Developer/CommandLineTools && xcode-select --install
# Optional: brew install make && gmake api && gmake config
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

TOOLS_BIN="${ROOT}/.tools/bin"
GOPATH_BIN="$(go env GOPATH)/bin"
export PATH="${TOOLS_BIN}:${GOPATH_BIN}:${PATH}"

if ! command -v protoc >/dev/null 2>&1; then
  echo "protoc not found. Install: brew install protobuf" >&2
  exit 1
fi

need_plugin() {
  command -v "$1" >/dev/null 2>&1
}
for p in protoc-gen-go protoc-gen-go-grpc protoc-gen-go-http protoc-gen-openapi; do
  if ! need_plugin "$p"; then
    cat >&2 <<EOF
Missing ${p} in PATH. Install protoc plugins (once), e.g.:

  export GOBIN=\$(go env GOPATH)/bin   # must be writable; or: export GOBIN=${ROOT}/.tools/bin
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
  go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

When GNU make works: make init
EOF
    exit 1
  fi
done

INTERNAL_PROTO_FILES=$(find internal -name '*.proto' 2>/dev/null || true)
API_PROTO_FILES=$(find api -name '*.proto' 2>/dev/null || true)

if [[ -n "${API_PROTO_FILES}" ]]; then
  echo "protoc api ..."
  # shellcheck disable=SC2086
  protoc --proto_path=./api \
    --proto_path=./third_party \
    --go_out=paths=source_relative:./api \
    --go-http_out=paths=source_relative:./api \
    --go-grpc_out=paths=source_relative:./api \
    --openapi_out=fq_schema_naming=true,default_response=false:. \
    ${API_PROTO_FILES}
fi

if [[ -n "${INTERNAL_PROTO_FILES}" ]]; then
  echo "protoc internal (config) ..."
  # shellcheck disable=SC2086
  protoc --proto_path=./internal \
    --proto_path=./third_party \
    --go_out=paths=source_relative:./internal \
    ${INTERNAL_PROTO_FILES}
fi

echo "go generate + tidy ..."
go generate ./...
go mod tidy

echo "Done."
