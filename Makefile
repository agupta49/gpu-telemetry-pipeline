APP_NAME=gpu-telemetry
IMG_PREFIX ?= localhost
NAMESPACE=gpu-telemetry
CLUSTER_NAME=gpu-telemetry

.PHONY: test cover cover-html swagger build-binaries docker-build kind-create kind-delete kind-load helm-install helm-uninstall lint clean logs-streamer

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

build-binaries:
	go build -o bin/streamer ./cmd/streamer
	go build -o bin/collector ./cmd/collector
	go build -o bin/mq ./cmd/mq
	go build -o bin/api-gateway ./cmd/api-gateway

docker-build: # For speed: DOCKER_DEFAULT_PLATFORM=linux/amd64 DOCKER_BUILDKIT=1 make docker-build -j4
	docker build -f deploy/docker/Dockerfile.streamer -t $(IMG_PREFIX)/streamer:latest .
	docker build -f deploy/docker/Dockerfile.collector -t $(IMG_PREFIX)/collector:latest .
	docker build -f deploy/docker/Dockerfile.mq -t $(IMG_PREFIX)/mq:latest .
	docker build -f deploy/docker/Dockerfile.api-gateway -t $(IMG_PREFIX)/api-gateway:latest

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
