#!/bin/bash
set -e
echo "=== Building Linux binaries for Docker/k8s ==="
echo "Host OS: $(go env GOOS)/$(go env GOARCH)"
echo "Target: linux/amd64"
mkdir -p bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/streamer ./cmd/streamer
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/collector ./cmd/collector
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/mq ./cmd/mq
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api-gateway ./cmd/api-gateway
echo "Done. Verify they're Linux binaries:"
file bin/streamer
ls -lh bin/
