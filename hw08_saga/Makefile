BIN_CMD_ORDER := "./bin/order_service"
BIN_CMD_PAY := "./bin/pay_service"
BIN_CMD_STOCK := "./bin/stock_service"
BIN_CMD_DELIVERY := "./bin/delivery_service"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CMD_ORDER) -ldflags "$(LDFLAGS)" ./cmd/order_service
	go build -v -o $(BIN_CMD_PAY) -ldflags "$(LDFLAGS)" ./cmd/pay_service  
	go build -v -o $(BIN_CMD_STOCK) -ldflags "$(LDFLAGS)" ./cmd/stock_service  
	go build -v -o $(BIN_CMD_DELIVERY) -ldflags "$(LDFLAGS)" ./cmd/delivery_service

run: build
	$(BIN_CMD_ORDER) -config ./configs/order_config.toml

run-pay: build
	$(BIN_CMD_PAY) -config ./configs/pay_config.toml

run-stock: build
	$(BIN_CMD_STOCK) -config ./configs/stock_config.toml

run-delivery: build
	$(BIN_CMD_DELIVERY) -config ./configs/delivery_config.toml

dbuild:
	docker compose -f ./deployments/docker-compose.yaml build --no-cache

up:
	docker compose -f ./deployments/docker-compose.yaml up

down:
	docker compose -f ./deployments/docker-compose.yaml down --rmi all

version: build
	$(BIN_CMD_ORDER) version
	$(BIN_CMD_PAY) version
	$(BIN_CMD_STOCK) version
	$(BIN_CMD_DELIVERY) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
