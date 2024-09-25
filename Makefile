GOCYCLO_PATH=$(shell go env GOPATH)/bin/gocyclo
GOCOGNIT_PATH=$(shell go env GOPATH)/bin/gocognit
GODEPGRAPH_PATH=$(shell go env GOPATH)/bin/godepgraph
GOOSEE_PATH=$(shell go env GOPATH)/bin/goose

MIGRATIONS_PATH=./migrations
POSTGRESQL_URI="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

THRESHOLD=5

BIN_DIR=bin
APP_NAME=manager

APP_PATH_SRC=cmd/$(APP_NAME)/main.go
APP_PATH_BIN=$(BIN_DIR)/$(APP_NAME)

.PHONY: all mkdir-bin run build tidy clean gocyclo gocognit test coverage
.PHONY: unit-test integration-test e2e-test benchmark

all: build

run: build
	./$(APP_PATH_BIN)

unit-test:
	@echo "Unit Tests:"
	@go test ./internal/usecase/ -coverprofile=coverage_usecase.out

integration-test:
	@echo "Integration Tests:"
	@go test -coverpkg=./internal/storage/storage_json -coverprofile=coverage_storage.out ./tests/integration_test.go

e2e-test: build
	@echo "E2E Tests:"
	@go test tests/e2e_test.go

test: unit-test integration-test e2e-test
	@echo "mode: set" > coverage.out
	@tail -n +2 coverage_usecase.out >> coverage.out
	@tail -n +2 coverage_storage.out >> coverage.out
	@rm coverage_usecase.out coverage_storage.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html 

benchmark:
	@go test -bench=. -benchtime=10x  -benchmem ./benchmark/storage_test.go

build: dependancy-install mkdir-bin $(APP_PATH_BIN) gocyclo gocognit

mkdir-bin:
	@mkdir -p bin

$(APP_PATH_BIN): mkdir-bin
	go build -o $(APP_PATH_BIN) $(APP_PATH_SRC)

dependancy-update:
	@go get -u

dependancy-install:
	@go get ./internal/... ./benchmark/... ./tests/...

tidy:
	@go mod tidy

gocyclo-install:
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gocognit-install:
	@go install github.com/uudashr/gocognit/cmd/gocognit@latest

gocyclo: gocyclo-install
	$(GOCYCLO_PATH) -over $(THRESHOLD) -ignore "_mock|_test" internal

gocognit: gocognit-install
	$(GOCOGNIT_PATH) -over $(THRESHOLD) -ignore "_mock|_test" internal

depgraph-install:
	@go install github.com/kisielk/godepgraph@latest

depgraph-build:
	$(GODEPGRAPH_PATH) -s $(APP_PATH_SRC) | dot -Tpng -o godepgraph.png

depgraph: depgraph-install depgraph-build

compose-up:
	docker-compose up -d postgres

compose-down:
	docker-compose down

compose-stop:
	docker-compose stop postgres

compose-start:
	docker-compose start postgres

compose-ps:
	docker-compose ps postgres

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_URI) create rename_me sql

goose-up:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_URI) up

goose-down:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_URI) down

goose-status:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_URI) status

.PHONY: depgraph compose-up compose-down compose-stop compose-start goose-install goose-add goose-up goose-status goose-down

clean:
	rm -rf $(BIN_DIR) godepgraph.png