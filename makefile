ACCRUAL_PORT := $(shell random unused-port)
GOPHERMART_PORT := $(shell random unused-port)
DATABASE_ACCRUAL := postgres://postgres:postgres@localhost:5432/accrual?sslmode=disable
DATABASE_GOPHERMART := postgres://postgres:postgres@localhost:5433/gophermart?sslmode=disable

.DEFAULT_GOAL := all

.PHONY: all
all: test lint autotest

.PHONY: generate
generate:
	@go generate ./...

.PHONY: up
up:
	@gophermartLogLevel=info 
		gophermartRunAddress="0.0.0.0:8080" \
		gophermartDBURI="postgresql://postgres:postgres@gophermartpostgres:5432/gophermart?sslmode=disable" \
		accrualAddress="localhost:8081" \
		secretKeyPath="./secret_key.txt" \
		tokenTTL=10m \
		accrualLogLevel=info \
		accrualRunAddress="0.0.0.0:8081" \
		accrualDBURI="postgresql://postgres:postgres@accrualpostgres:5432/accrual?sslmode=disable" \
		docker-compose -f ./scripts/docker-compose.yaml up -d

.PHONY: down
down:
	@docker-compose -f ./scripts/docker-compose.yaml down

.PHONY: lint
lint:
	@go vet -vettool=$(shell which statictest) ./...

.PHONY: test
test:
	@go test -short -race -timeout=30s -count=1 -cover ./...

.PHONY: build
build:
	@go build -buildvcs=false -o ./cmd/accrual/accrual ./cmd/accrual
	@go build -buildvcs=false -o ./cmd/gophermart/gophermart ./cmd/gophermart

.PHONY: autotest
autotest: build
	@gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=$(GOPHERMART_PORT) \
		-gophermart-database-uri=$(DATABASE_GOPHERMART) \
		-accrual-binary-path=cmd/accrual/accrual \
		-accrual-host=localhost \
		-accrual-port=$(ACCRUAL_PORT) \
		-accrual-database-uri=$(DATABASE_ACCRUAL)
