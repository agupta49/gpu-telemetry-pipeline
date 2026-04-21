#!/bin/bash
set -e
echo "=== Building binaries on host ==="
mkdir -p bin
echo "Building streamer..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/streamer ./cmd/streamer
echo "Building collector..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/collector ./cmd/collector
echo "Building mq..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/mq ./cmd/mq
echo "Building api-gateway..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api-gateway ./cmd/api-gateway
echo "Done. ls -lh bin/"
ls -lh bin/
