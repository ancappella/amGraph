#!/usr/bin/env bash
# Same as: make api && make config && go generate ./...
# Uses Homebrew GNU make (gmake) so it works when Apple /usr/bin/make fails with xcrun/CLT errors.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

TOOLS="${ROOT}/.tools/bin"
mkdir -p "${TOOLS}"
export GOBIN="${TOOLS}"
export PATH="${TOOLS}:$(go env GOPATH)/bin:${PATH}"

if ! command -v protoc-gen-go-http >/dev/null 2>&1; then
  echo "Installing protoc plugins into ${TOOLS} ..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
  go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
fi

if command -v gmake >/dev/null 2>&1; then
  MAKE=gmake
elif command -v gnumake >/dev/null 2>&1; then
  MAKE=gnumake
else
  MAKE=make
fi

${MAKE} api
${MAKE} config
go generate ./...
go mod tidy
echo "codegen OK"
