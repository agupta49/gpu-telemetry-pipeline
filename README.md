# **Elastic GPU Telemetry Pipeline**

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Ingests DCGM CSV streams, routes via custom gRPC message queue, persists to Postgres/TimescaleDB, and exposes REST APIs.

### **Prerequisites**
**Option 1: Auto-install script for Amazon Linux 2023 / RHEL**
```bash
bash scripts/install-prereqs.sh
```

**Option 2: Manual install**
1. **Go 1.22+**: `sudo dnf install -y golang`
2. **Docker**: `sudo dnf install -y docker && sudo systemctl start docker && sudo systemctl enable docker`
3. **kubectl**: `curl -LO https://dl.k8s.io/release/v1.30.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/`
4. **Helm**: `curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash`
5. **kind**: `curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64 && chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind`

### **Quick Start**
```bash
git clone https://github.com/agupta49/gpu-telemetry-pipeline.git
cd gpu-telemetry-pipeline
go mod tidy

# Docker perms
sudo usermod -aG docker $USER && newgrp docker

# Place your DCGM CSV in data/ if you have one
cp /path/to/dcgm_metrics_20250718_134233.csv data/

# Tests >80%
make cover

# Deploy to kind
make kind-create
DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 make docker-build -j4 kind-load
make helm-install

# Verify
kubectl get pods -n gpu-telemetry -w
curl http://localhost:30080/api/v1/gpus
```

### **Make Targets**
| Command | Description |
| --- | --- |
| `make cover` | Run tests, fail if <80% coverage |
| `make kind-create` | Create kind cluster |
| `make docker-build kind-load` | Build and load images. Use `DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 -j4` for speed |
| `make helm-install` | Deploy Helm chart |
| `make logs-streamer` | Tail streamer logs |

### **Coverage**
`make cover` only tests `internal/*` packages. `cmd/*` and `pkg/pb` excluded as wrappers/generated.

### **Troubleshooting**

**Docker build stuck at `go build` for hours**: This happens if Docker tries to emulate ARM on x86 via QEMU. Fix:

1. **Kill the stuck build**: 
   ```bash
   docker buildx prune -a -f
   docker system prune -a
   ```
2. **Force native platform**: 
   ```bash
   DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 make docker-build -j4
   ```
3. **Check your arch**: `uname -m`. If `x86_64`, the `--platform=linux/amd64` flag in Dockerfile prevents emulation.
4. **Test native build**: `go build ./cmd/streamer` should take <10s. If slow, your machine is the issue.

**Docker permission denied**: Run `newgrp docker` or logout/login after `usermod -aG docker $USER`

**Kubernetes cluster unreachable**: Run `make kind-create` to create the cluster
