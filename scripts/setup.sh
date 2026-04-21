#!/bin/bash
set -e

echo "=== GPU Telemetry Pipeline Setup ==="

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "Go not found. Install: sudo dnf install -y golang"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker not found. Install: sudo dnf install -y docker"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "kubectl not found. See README for install"; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "helm not found. See README for install"; exit 1; }
command -v kind >/dev/null 2>&1 || { echo "kind not found. See README for install"; exit 1; }

echo "✓ All prerequisites found"

# Docker group
if ! groups | grep -q docker; then
    echo "Adding user to docker group..."
    sudo usermod -aG docker $USER
    echo "Run 'newgrp docker' or logout/login to apply"
fi

# Start docker
sudo systemctl start docker
sudo systemctl enable docker

echo "✓ Docker running"

# Go mod
go mod tidy
go mod verify

echo "✓ Go modules ready"
echo ""
echo "Next steps:"
echo "1. make cover          # Run tests"
echo "2. make kind-create    # Create cluster"
echo "3. make docker-build kind-load"
echo "4. make helm-install"
