# **Elastic GPU Telemetry Pipeline**

Repo: https://github.com/agupta49/gpu-telemetry-pipeline.git

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Ingests DCGM CSV streams, routes via custom gRPC message queue, persists to Postgres/TimescaleDB, and exposes REST APIs.

### **Prerequisites**

1. **Go 1.22+**: `sudo dnf install -y golang`
2. **Docker**: `sudo dnf install -y docker && sudo systemctl start docker && sudo systemctl enable docker`
3. **kubectl**: `curl -LO https://dl.k8s.io/release/v1.30.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/`
4. **Helm**: `curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash`
5. **kind**: `curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64 && chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind`

### **Quick Start - Local kind Cluster**

```bash
# 1. Clone and setup
git clone https://github.com/agupta49/gpu-telemetry-pipeline.git
cd gpu-telemetry-pipeline
go mod tidy

# 2. Add user to docker group (logout/login or newgrp after)
sudo usermod -aG docker $USER
newgrp docker

# 3. Run tests - must show >80%
make cover

# 4. Create kind cluster
make kind-create

# 5. Build images and load into kind
make docker-build kind-load

# 6. Deploy to kind
make helm-install

# 7. Verify
kubectl get pods -n gpu-telemetry -w
curl http://localhost:30080/api/v1/gpus
```

### **Make Targets**

| Command | Description |
| --- | --- |
| `make cover` | Run unit tests with coverage check. Fails if <80% |
| `make cover-html` | Generate HTML coverage report |
| `make kind-create` | Create kind cluster named gpu-telemetry |
| `make kind-delete` | Delete kind cluster |
| `make docker-build` | Build all Docker images |
| `make kind-load` | Load built images into kind cluster |
| `make helm-install` | Deploy Helm chart with wait |
| `make helm-uninstall` | Remove Helm release and namespace |
| `make logs-streamer` | Tail streamer logs |
| `make clean` | Remove build artifacts |

### **Database Auto-Setup**
Postgres + TimescaleDB installed automatically via Helm. ConfigMap `gpu-telemetry-init-db` creates extension, `telemetry` hypertable, and indexes. `collector` waits for table before starting.

### **Access API**
After `make helm-install`, the API is exposed on NodePort 30080:
```bash
curl http://localhost:30080/api/v1/gpus
curl http://localhost:30080/healthz
```

### **Troubleshooting**

**Docker permission denied**: Run `newgrp docker` or logout/login after `usermod -aG docker $USER`

**Kubernetes cluster unreachable**: Run `make kind-create` to create the cluster

**Image not found**: Make sure you ran `make docker-build kind-load` before `make helm-install`

**Coverage <80%**: Run `make cover-html` and open `coverage.html` to see uncovered lines

### **Coverage**
Run `make cover` to test only internal packages. `cmd/*` and `pkg/pb` are excluded from coverage as they are wrappers or generated code.
