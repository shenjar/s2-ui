#!/bin/sh
# Build a static, drop-in `sui` binary for TESTING on a Linux server.
#
# It mirrors the official release build (static musl) but DROPS the naive /
# cronet outbound, which avoids a very heavy toolchain. The result is a single
# static file that runs on any Linux of the matching arch. naive outbound will
# not work in this test binary — irrelevant for testing the stats chart and the
# ACME auto-SSL feature. (If you actually use a naive outbound, tell me and I'll
# give you the full build instead.)
#
# Usage (from the repo root, AFTER building the frontend into web/html):
#   cd frontend && npm install && npm run build && cd ..
#   rm -rf web/html && mkdir -p web/html && cp -r frontend/dist/* web/html/
#   docker run --rm -v "$PWD":/app -w /app golang:1.26-alpine sh buildtest.sh
#
# Other arch (e.g. arm64 server):
#   docker run --rm -e GOARCH=arm64 -v "$PWD":/app -w /app golang:1.26-alpine sh buildtest.sh
set -e

apk add --no-cache gcc musl-dev git >/dev/null

: "${GOARCH:=amd64}"
export CGO_ENABLED=1 GOOS=linux GOARCH

go build -ldflags="-w -s -checklinkname=0 -linkmode external -extldflags '-static'" \
  -tags "with_quic,with_grpc,with_utls,with_acme,with_gvisor,badlinkname,tfogo_checklinkname0,with_tailscale" \
  -o sui main.go

echo "Built ./sui (GOARCH=$GOARCH). Copy it to /usr/local/s-ui/sui on the server and restart s-ui."
