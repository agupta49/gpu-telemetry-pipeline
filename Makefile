APP_NAME=gpu-telemetry
IMG_PREFIX ?= ghcr.io/agupta49
NAMESPACE=gpu-telemetry

.PHONY: test cover cover-html swagger build-binaries docker-build kind-load helm-install helm-uninstall lint clean

test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

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

docker-build:
	docker build -f deploy/docker/Dockerfile.streamer -t $(IMG_PREFIX)/streamer:latest .
	docker build -f deploy/docker/Dockerfile.collector -t $(IMG_PREFIX)/collector:latest .
	docker build -f deploy/docker/Dockerfile.mq -t $(IMG_PREFIX)/mq:latest .
	docker build -f deploy/docker/Dockerfile.api-gateway -t $(IMG_PREFIX)/api-gateway:latest

kind-load:
	kind load docker-image $(IMG_PREFIX)/streamer:latest --name gpu-telemetry
	kind load docker-image $(IMG_PREFIX)/collector:latest --name gpu-telemetry
	kind load docker-image $(IMG_PREFIX)/mq:latest --name gpu-telemetry
	kind load docker-image $(IMG_PREFIX)/api-gateway:latest --name gpu-telemetry

helm-install:
	helm upgrade --install $(APP_NAME) ./charts/gpu-telemetry \
		--namespace $(NAMESPACE) --create-namespace \
		--set global.imagePrefix=$(IMG_PREFIX) \
		--wait --timeout 5m

helm-uninstall:
	helm uninstall $(APP_NAME) -n $(NAMESPACE) || true
	kubectl delete ns $(NAMESPACE) || true

clean:
	rm -rf bin/ coverage.out coverage.html
