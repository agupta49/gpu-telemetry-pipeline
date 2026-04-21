APP_NAME=gpu-telemetry
IMG_PREFIX ?= localhost
NAMESPACE=gpu-telemetry
CLUSTER_NAME=gpu-telemetry

TARGET_GOOS := linux
TARGET_GOARCH := amd64

.PHONY: test cover cover-html swagger proto build-binaries docker-build kind-create kind-delete kind-load helm-install helm-uninstall lint clean logs-streamer build-all

PKG_LIST := $(shell go list ./... | grep -v /cmd/ | grep -v /pkg/pb)

test:
	go test $(PKG_LIST) -coverprofile=coverage.out -covermode=atomic

cover: test
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@go tool cover -func=coverage.out | grep total | awk '{exit ($$3+0)<80.0?1:0}'

cover-html: test
	go tool cover -html=coverage.out -o coverage.html

swagger:
	swag init -g cmd/api-gateway/main.go -o api --parseDependency --parseInternal

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/pb/telemetry.proto

build-binaries:
	@echo "Building Linux binaries for Docker/k8s"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/streamer ./cmd/streamer
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/collector ./cmd/collector
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/mq ./cmd/mq
	CGO_ENABLED=0 GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) go build -ldflags="-s -w" -o bin/api-gateway ./cmd/api-gateway

docker-build:
	docker build -f deploy/docker/Dockerfile.streamer -t $(IMG_PREFIX)/streamer:latest .
	docker build -f deploy/docker/Dockerfile.collector -t $(IMG_PREFIX)/collector:latest .
	docker build -f deploy/docker/Dockerfile.mq -t $(IMG_PREFIX)/mq:latest .
	docker build -f deploy/docker/Dockerfile.api-gateway -t $(IMG_PREFIX)/api-gateway:latest

build-all: build-binaries docker-build

kind-create:
	kind create cluster --name $(CLUSTER_NAME) --config=- <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 30080
EOF

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
