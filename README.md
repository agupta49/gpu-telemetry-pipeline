# **Elastic GPU Telemetry Pipeline**

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Works on macOS, Linux, and EC2.

### **Why This Build Is Fast**

We cross-compile Go binaries on your host machine. Docker images just `COPY` the 5MB Linux binary. Works identically on macOS M1/M2/M3, Intel Mac, Linux, or Windows WSL.

### **Prerequisites**

**macOS:**
```bash
brew install go docker kind helm kubectl
```

**Linux/RHEL/EC2:**
```bash
bash scripts/install-prereqs.sh
newgrp docker
```

### **Quick Start - Works on macOS & Linux**
```bash
# 1. Extract and setup
unzip gpu-telemetry-pipeline.zip
cd gpu-telemetry-pipeline
go mod tidy

# 2. Add your DCGM CSV
cp /path/to/dcgm_metrics_20250718_134233.csv data/

# 3. Run tests - builds for your local OS
make cover

# 4. Build Linux binaries for Docker - cross-compiles on macOS automatically
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

### **Key Commands**

| Command | Description |
| --- | --- |
| `make build-binaries` | Cross-compile Linux/amd64 binaries. Works on macOS/Linux |
| `make build-binaries-local` | Build for your local OS to test natively |
| `make docker-build` | Build Docker images by copying Linux binaries |
| `make build-all` | Does both above |
| `make cover` | Run tests, fail if <80% coverage |

### **macOS Specific Notes**

1. **Docker Desktop required**: Install from docker.com. Start it before `make docker-build`
2. **Cross-compilation is automatic**: `make build-binaries` sets `GOOS=linux GOARCH=amd64` so binaries work in k8s even on M1/M2/M3 Macs
3. **Test locally**: Run `make build-binaries-local` to get `bin/streamer-darwin_arm64` you can run directly on Mac
4. **Performance**: First build compiles deps ~2-3 min. Second build uses cache ~3s

### **Troubleshooting**

**`exec format error` in pod**: You ran `docker-build` without `build-binaries` first. Docker copied wrong arch. Fix: `make clean && make build-binaries && make docker-build`

**Build slow on macOS**: First run compiles `grpc` + `golang.org/x/text` ~2 min. Second run uses `~/Library/Caches/go-build` ~3s. If still slow after second run, check `rm -rf vendor/`

**`docker: command not found` on Mac**: Install Docker Desktop and start it. Check `docker ps` works
