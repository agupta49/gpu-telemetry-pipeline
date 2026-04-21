#!/bin/bash
set -e

echo "=== Installing Prerequisites for GPU Telemetry Pipeline ==="
echo "OS: Amazon Linux 2023 / RHEL / CentOS"

# Detect if running as root
if [ "$EUID" -eq 0 ]; then
  SUDO=""
else
  SUDO="sudo"
fi

# 1. Update system
echo "[1/5] Updating system packages..."
$SUDO dnf update -y

# 2. Install Go 1.22+
echo "[2/5] Installing Go 1.22..."
if ! command -v go >/dev/null 2>&1; then
  $SUDO dnf install -y golang
  echo "Go installed: $(go version)"
else
  echo "Go already installed: $(go version)"
fi

# 3. Install Docker
echo "[3/5] Installing Docker..."
if ! command -v docker >/dev/null 2>&1; then
  $SUDO dnf install -y docker
  $SUDO systemctl start docker
  $SUDO systemctl enable docker
  $SUDO usermod -aG docker $USER
  echo "Docker installed. Run 'newgrp docker' or logout/login to use without sudo"
else
  echo "Docker already installed: $(docker --version)"
  $SUDO systemctl start docker || true
fi

# 4. Install kubectl
echo "[4/5] Installing kubectl..."
if ! command -v kubectl >/dev/null 2>&1; then
  KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
  curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl"
  chmod +x kubectl
  $SUDO mv kubectl /usr/local/bin/
  echo "kubectl installed: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
else
  echo "kubectl already installed: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
fi

# 5. Install Helm
echo "[5/5] Installing Helm..."
if ! command -v helm >/dev/null 2>&1; then
  curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
  chmod 700 get_helm.sh
  ./get_helm.sh
  rm get_helm.sh
  echo "Helm installed: $(helm version --short)"
else
  echo "Helm already installed: $(helm version --short)"
fi

# 6. Install kind
echo "[6/6] Installing kind..."
if ! command -v kind >/dev/null 2>&1; then
  KIND_VERSION="v0.23.0"
  curl -Lo ./kind "https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64"
  chmod +x ./kind
  $SUDO mv ./kind /usr/local/bin/kind
  echo "kind installed: $(kind version)"
else
  echo "kind already installed: $(kind version)"
fi

echo ""
echo "=== All Prerequisites Installed ==="
echo "Versions:"
echo "  Go:      $(go version 2>/dev/null || echo 'not found')"
echo "  Docker:  $(docker --version 2>/dev/null || echo 'not found')"
echo "  kubectl: $(kubectl version --client --short 2>/dev/null || echo 'not found')"
echo "  Helm:    $(helm version --short 2>/dev/null || echo 'not found')"
echo "  kind:    $(kind version 2>/dev/null || echo 'not found')"
echo ""
echo "IMPORTANT: If docker was just installed, run 'newgrp docker' or logout/login"
echo "Next: cd gpu-telemetry-pipeline && make cover"
