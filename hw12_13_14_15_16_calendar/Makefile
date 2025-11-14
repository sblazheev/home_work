BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: swag
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) serve --config "./configs/config.yaml"

run-grpc: build
	$(BIN) grpc --config "./configs/config.yaml"

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... -tags !integrations

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.6.1

lint: install-lint-deps
	golangci-lint run ./...

migrate: build
	$(BIN) migrate --config "./configs/config.yaml"

install-swag-deps:
	(which swag > /dev/null) || go install github.com/swaggo/swag/cmd/swag@latest

swag: install-swag-deps
	swag i -d internal/server/http/,internal/app/,internal/common/ -o internal/server/http/docs/ -g server.go

generate:
	rm -rf proto/pb
	mkdir -p proto/pb

	protoc --proto_path=internal/server/grpc/proto/ --go_out=internal/server/grpc/pb --go-grpc_out=internal/server/grpc/pb internal/server/grpc/proto/*.proto

.PHONY: build run build-img run-img version test lint
