# **Elastic GPU Telemetry Pipeline**

Repo: https://github.com/agupta49/gpu-telemetry-pipeline.git

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Ingests DCGM CSV streams, routes via custom gRPC message queue, persists to Postgres/TimescaleDB, and exposes REST APIs.

### **Module Path**
`module github.com/agupta49/gpu-telemetry-pipeline`

### **Build & Deploy to kind**
```bash
go mod tidy
IMG_PREFIX=localhost make docker-build kind-load
make helm-install
```

### **Database Auto-Setup**
Postgres + TimescaleDB installed automatically via Helm. ConfigMap `gpu-telemetry-init-db` creates extension, `telemetry` hypertable, and indexes. `collector` waits for table before starting.

### **Verify**
```bash
curl http://localhost:8080/api/v1/gpus
make cover # must show >80%
```
