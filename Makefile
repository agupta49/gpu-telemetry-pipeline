APP_NAME=gpu-telemetry
IMG_PREFIX ?= localhost
NAMESPACE=gpu-telemetry
CLUSTER_NAME=gpu-telemetry

# Detect OS for local binary names
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LOCAL_BIN_SUFFIX := $(GOOS)_$(GOARCH)

# Target OS/ARCH for Docker images - always linux/amd64 for k8s
TARGET_GOOS := linux
TARGET_GOARCH := amd64

.PHONY: test cover cover-html swagger build-binaries build-binaries-local docker-build kind-create kind-delete kind-load helm-install helm-uninstall lint clean logs-streamer build-all

PKG_LIST := $(shell go list ./... | grep -v /cmd/ | grep -v /pkg/pb)

test:
	go test $(PKG_LIST) -coverprofile=coverage.out -covermode=atomic

cover: test
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@go tool cover -func=coverage.out | grep total | awk '{exit ($$3+0)<80.0?1:0}'

cover-html: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in browser"

swagger:
	swag init -g cmd/api-gateway/main.go -o api --parseDependency --parseInternal

build-binaries: # Cross-compile for Linux - works on macOS/Linux/Windows
	@echo "Building Linux binaries for Docker/k8s on $(GOOS)/$(GOARCH) host"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/streamer ./cmd/streamer
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/collector ./cmd/collector
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/mq ./cmd/mq
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/api-gateway ./cmd/api-gateway
	@echo "Linux binaries built: file bin/streamer"

build-binaries-local: # Build for local OS for testing
	@echo "Building local binaries for $(GOOS)/$(GOARCH)"
	@mkdir -p bin
	go build -o bin/streamer-$(LOCAL_BIN_SUFFIX) ./cmd/streamer
	go build -o bin/collector-$(LOCAL_BIN_SUFFIX) ./cmd/collector
	go build -o bin/mq-$(LOCAL_BIN_SUFFIX) ./cmd/mq
	go build -o bin/api-gateway-$(LOCAL_BIN_SUFFIX) ./cmd/api-gateway
	@echo "Local binaries built: ls bin/*$(LOCAL_BIN_SUFFIX)"

docker-build: # Builds are instant - just copies pre-built Linux binaries
	docker build -f deploy/docker/Dockerfile.streamer -t $(IMG_PREFIX)/streamer:latest .
	docker build -f deploy/docker/Dockerfile.collector -t $(IMG_PREFIX)/collector:latest .
	docker build -f deploy/docker/Dockerfile.mq -t $(IMG_PREFIX)/mq:latest .
	docker build -f deploy/docker/Dockerfile.api-gateway -t $(IMG_PREFIX)/api-gateway:latest

build-all: build-binaries docker-build

kind-create:
	kind create cluster --name $(CLUSTER_NAME)
	kubectl cluster-info --context kind-$(CLUSTER_NAME)

kind-delete:
	kind delete cluster --name $(CLUSTER_NAME)

kind-load:
	kind load docker-image $(IMG_PREFIX)/streamer:latest --name $(CLUSTER_NAME)
	kind load docker-image $(IMG_PREFIX)/collector:latest --name $(CLUSTER_NAME)
	kind load docker-image $(IMG_PREFIX)/mq:latest --name $(CLUSTER_NAME)
	kind load docker-image $(IMG_PREFIX)/api-gateway:latest --name $(CLUSTER_NAME)

helm-install:
	helm upgrade --install $(APP_NAME) ./charts/gpu-telemetry \
		--namespace $(NAMESPACE) --create-namespace \
		--set global.imagePrefix=$(IMG_PREFIX) \
		--wait --timeout 5m

helm-uninstall:
	helm uninstall $(APP_NAME) -n $(NAMESPACE) || true
	kubectl delete ns $(NAMESPACE) || true

logs-streamer:
	kubectl logs -n $(NAMESPACE) -l app.kubernetes.io/name=streamer --tail=50 -f

clean:
	rm -rf bin/ coverage.out coverage.html
