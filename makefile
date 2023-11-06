ACCRUAL_PORT := $(shell random unused-port)
GOPHERMART_PORT := $(shell random unused-port)
DATABASE_ACCRUAL := postgres://postgres:postgres@localhost:5432/accrual?sslmode=disable
DATABASE_GOPHERMART := postgres://postgres:postgres@localhost:5433/gophermart?sslmode=disable

.DEFAULT_GOAL := all

.PHONY: all
all: test lint autotest

.PHONY: up
up:
	@docker-compose -f ./scripts/docker-compose.yaml up -d

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
	@go build -buildvcs=false -o ./cmd/gophermart/server ./cmd/gophermart

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
