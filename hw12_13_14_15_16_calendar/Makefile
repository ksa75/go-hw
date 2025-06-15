BIN := "./bin/calendar"
BIN_SENDER := "./bin/calendar_sender"
BIN_SCHEDULER := "./bin/calendar_scheduler"
DOCKER_IMG="calendar:develop"
DOCKER_IMG_SENDER="calendar_sender:develop"
DOCKER_IMG_SCHEDULER="calendar_scheduler:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up:
	cd deployments
	docker-compose up -d

down:
	cd deployments
	docker-compose down

generate: install-protoc-deps
	rm -rf api/calendarpb
	protoc   \
	--proto_path=.   \
	--proto_path=/usr/include   \
	--go_out=api   \
	--go-grpc_out=api   \
	api/EventService.proto

install-protoc-deps:
	sudo apt install protobuf-compiler -y
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/Dockerfile.sender .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/Dockerfile.scheduler .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version
	$(BIN_SCHEDULER) version
	$(BIN_SENDER) version

test:
	go test -race ./internal/logger/... ./internal/storage/memory/... ./internal/server/http/... ./internal/server/grpc/... ./internal/scheduler/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

lint: install-lint-deps
	golangci-lint run ./... --fix

.PHONY: build run build-img run-img version test lint
