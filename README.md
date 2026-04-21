# **Elastic GPU Telemetry Pipeline**

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Ingests DCGM CSV streams, routes via custom gRPC message queue, persists to Postgres/TimescaleDB, and exposes REST APIs.

### **Why This Build Is Fast**

We compile Go binaries on your host machine where Go build cache lives. Docker images just `COPY` the 5MB binary. This is 60x faster than compiling inside Docker.

### **Prerequisites**
```bash
bash scripts/install-prereqs.sh
newgrp docker
```

### **Quick Start - New Workflow**
```bash
# 1. Extract and setup
unzip gpu-telemetry-pipeline.zip
cd gpu-telemetry-pipeline
go mod tidy

# 2. Add your DCGM CSV
cp /path/to/dcgm_metrics_20250718_134233.csv data/

# 3. Run tests
make cover

# 4. Build binaries on host - uses Go cache, ~10s total
make build-binaries

# 5. Build Docker images - instant, just copies binaries
make docker-build

# 6. Create cluster and load images
make kind-create
make kind-load

# 7. Deploy
make helm-install

# 8. Verify
kubectl get pods -n gpu-telemetry -w
curl http://localhost:30080/healthz
```

### **One-Liner**
```bash
make build-binaries && make docker-build && make kind-create && make kind-load && make helm-install
```

### **Make Targets**
| Command | Description |
| --- | --- |
| `make build-binaries` | Compile all 4 Go services on host using cache. ~10s |
| `make docker-build` | Build Docker images by copying binaries. ~2s |
| `make build-all` | Does both above |
| `make cover` | Run tests, fail if <80% coverage |
| `make kind-create` | Create kind cluster |
| `make kind-load` | Load images into kind |
| `make helm-install` | Deploy Helm chart |
| `make logs-streamer` | Tail streamer logs |

### **Troubleshooting**

**`go build` slow on host**: First run compiles deps, ~30s. Second run uses cache, ~2s. If still slow, your machine is underpowered.

**`exec format error` in pod**: You built binaries for wrong OS. Run `file bin/streamer` - should say `ELF 64-bit LSB executable, x86-64`. If it says `Mach-O` or `PE32`, you didn't set `GOOS=linux GOARCH=amd64`.

**Docker build slow**: You're using old Dockerfiles that compile inside Docker. Use `make build-binaries` first, then `make docker-build`. The new Dockerfiles only `COPY bin/streamer`.
