# **Elastic GPU Telemetry Pipeline with Custom Message Queue**

Elastic, scalable telemetry pipeline for AI clusters built in Golang. Ingests DCGM CSV streams, routes via custom gRPC message queue, persists to Postgres/TimescaleDB, and exposes REST APIs.

### **System Architecture**

```
[Streamer] -> gRPC -> [MQ] -> gRPC stream -> [Collector] -> Postgres/TimescaleDB
                                        ^
                                        |
                                  [API Gateway] -> REST /api/v1/gpus
```

**Components:**
1. **Telemetry Streamer**: Reads DCGM CSV, loops it, publishes each row to custom MQ. Horizontally scalable
2. **Messaging Queue**: Custom gRPC pub/sub service. Not Kafka/RabbitMQ. Designed for <10 instances
3. **Telemetry Collector**: Subscribes to MQ, parses telemetry, writes to TimescaleDB. Horizontally scalable  
4. **API Gateway**: Gin REST API with Swagger. Queries Postgres for GPU list + telemetry

**Design Considerations:**
- **Scale**: Streamer/Collector scale independently via K8s HPA. MQ is single service for this exercise but gRPC streaming allows fan-out
- **Performance**: Binary protobuf over gRPC, batch inserts in collector, TimescaleDB hypertables
- **Availability**: Health checks, retries with backoff in streamer. CSV loop provides continuous data if source stalls
- **Timestamp**: Per requirement, we use processing time as telemetry timestamp

### **API Endpoints**

| Method | Path | Description |
| --- | --- | --- |
| GET | `/healthz` | Liveness probe |
| GET | `/api/v1/gpus` | List all GPU IDs with telemetry |
| GET | `/api/v1/gpus/{id}/telemetry` | Get telemetry for GPU, ordered by time |
| GET | `/api/v1/gpus/{id}/telemetry?start_time=...&end_time=...` | Filter by time window, inclusive |

Time format: RFC3339 `2006-01-02T15:04:05Z`

### **Build and Package**

**Prerequisites macOS:** `brew install go docker kind helm kubectl protobuf`  
**Prerequisites Linux:** `bash scripts/install-prereqs.sh`

```bash
# 1. Generate protobuf
make proto

# 2. Build Linux binaries - works on macOS/Linux
make build-binaries

# 3. Build Docker images
make docker-build

# 4. Run tests with coverage gate
make cover  # fails if <80%

# 5. Generate OpenAPI spec
make swagger  # outputs to api/
```

### **Installation Workflow**

```bash
# 1. Create cluster with port mapping for NodePort
make kind-create

# 2. Load images
make kind-load

# 3. Deploy Helm chart
make helm-install

# 4. Verify
kubectl get pods -n gpu-telemetry -w
curl http://localhost:30080/healthz
curl http://localhost:30080/api/v1/gpus
```

### **Sample User Workflow**

```bash
# 1. Add DCGM CSV
cp dcgm_metrics_20250718_134233.csv data/

# 2. Deploy stack
make build-all && make kind-create && make kind-load && make helm-install

# 3. Stream starts automatically, loops CSV

# 4. Query data after 30s
curl "http://localhost:30080/api/v1/gpus"
curl "http://localhost:30080/api/v1/gpus/0/telemetry?start_time=2025-07-18T13:42:00Z"

# 5. Scale collectors
kubectl scale deployment gpu-telemetry-collector -n gpu-telemetry --replicas=3
```

### **Document Your Use of AI**

**Project Bootstrap**: Used AI to generate initial repo structure, Makefile, Helm charts. Prompt: "Create Go microservice monorepo with 4 services, Helm chart, Makefile with test/coverage/swagger targets." AI output required manual fixes for Helm template syntax and import paths.

**Code Bootstrap**: Used AI for boilerplate gRPC server/client, Gin handlers, Postgres connection. Prompt: "Go gRPC pubsub service with Publish and Subscribe streams." AI generated correct protobuf but missed error handling on stream sends - added manually.

**Unit Tests**: Used AI to generate table tests for db.go and handlers.go. Prompt: "Go table-driven tests for Config.Validate covering empty fields." AI output was 90% correct, manually added DB mock for repo tests.

**Build Env**: Used AI to write Dockerfile with multi-stage build. Hit BuildKit/QEMU slowness. Manual intervention: switched to pre-building binaries on host to use Go cache. This was the key fix.

**Where AI Fell Short**:
1. Helm templates: AI used wrong `{{ .Release.Name }}` syntax, fixed manually
2. Protobuf streaming: AI didn't handle context cancellation, added manually
3. Cross-compilation: AI suggested `GOOS=linux` but didn't explain vendor/ cache busting. Manual debug via `go build -x`
4. Test coverage: AI wrote tests but didn't hit 80% due to missing error paths. Added error case tests manually

**Manual Interventions**: Architecture decisions, performance tuning, debugging vendor cache issues, Kubernetes port mapping for NodePort access.
