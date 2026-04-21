#!/bin/bash
set -e
echo "=== GPU Telemetry Pipeline Setup ==="
command -v go >/dev/null 2>&1 || { echo "Go not found"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker not found"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "kubectl not found"; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "helm not found"; exit 1; }
command -v kind >/dev/null 2>&1 || { echo "kind not found"; exit 1; }
echo "✓ All prerequisites found"
go mod tidy
echo "Ready. Run: make cover && make kind-create && DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 make docker-build -j4 kind-load && make helm-install"
