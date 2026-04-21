#!/bin/bash
set -e
echo "=== Installing Prerequisites ==="
SUDO="sudo"; [ "$EUID" -eq 0 ] && SUDO=""
$SUDO dnf update -y
$SUDO dnf install -y golang docker
$SUDO systemctl start docker && $SUDO systemctl enable docker
$SUDO usermod -aG docker $USER
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && $SUDO mv kubectl /usr/local/bin/
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
curl -Lo ./kind "https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64"
chmod +x ./kind && $SUDO mv ./kind /usr/local/bin/kind
echo "Done. Run 'newgrp docker' then 'make build-binaries'"
