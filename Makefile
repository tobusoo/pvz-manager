include .env

PROTO_GENERATE_PATH=$(CURDIR)/pkg

GOCYCLO_PATH=$(shell go env GOPATH)/bin/gocyclo
GOCOGNIT_PATH=$(shell go env GOPATH)/bin/gocognit
GODEPGRAPH_PATH=$(shell go env GOPATH)/bin/godepgraph
GOOSEE_PATH=$(shell go env GOPATH)/bin/goose

MIGRATIONS_PATH=./migrations

THRESHOLD=5

BIN_DIR=$(CURDIR)/bin
APP_NAME=manager

APP_PATH_SRC=cmd/$(APP_NAME)/main.go
APP_PATH_BIN=$(BIN_DIR)/$(APP_NAME)

.PHONY: all mkdir-bin run build tidy clean gocyclo gocognit test coverage
.PHONY: unit-test integration-test integration-test-db e2e-test benchmark

all: bin-deps generate build run

run: build
	$(APP_PATH_BIN)

unit-test:
	@echo "Unit Tests:"
	@go test ./internal/usecase/ -coverprofile=coverage_usecase.out

integration-test-db:
	docker-compose up -d postgres_test
	@echo "Sleeping 4 seconds for postgreSQL preparation"
	@sleep 4
	@$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_TEST_DSN) up
	@POSTGRESQL_TEST_DSN=${POSTGRESQL_TEST_DSN} go test -v -coverpkg=./internal/storage/postgres \
		-coverprofile=coverage_storage_postgres.out \
		./tests/integration/storage_db/integration_test.go
	docker-compose down postgres_test

integration-test:
	@echo "Integration Tests:"
	@go test -coverpkg=./internal/storage/storage_json -coverprofile=coverage_storage.out \
		./tests/integration/storage_json/integration_test.go

e2e-test: build
	@echo "E2E Tests:"
	@go test tests/e2e_test.go

test: unit-test integration-test integration-test-db e2e-test
	@echo "mode: set" > coverage.out
	@tail -n +2 coverage_usecase.out >> coverage.out
	@tail -n +2 coverage_storage.out >> coverage.out
	@tail -n +2 coverage_storage_postgres.out >> coverage.out
	@rm coverage_usecase.out coverage_storage.out coverage_storage_postgres.out

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
	@go get ./internal/... ./benchmark/... ./tests/... ./scripts/... 

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
	docker-compose down postgres

compose-stop:
	docker-compose stop postgres

compose-start:
	docker-compose start postgres

compose-ps:
	docker-compose ps postgres

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN) create rename_me sql

goose-up:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN) up

goose-down:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN) down

goose-status:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN) status

squawk-install:
	npm install -g squawk-cli

squawk:
	squawk ./migrations/* --exclude=ban-drop-table

.PHONY: depgraph compose-up compose-down compose-stop compose-start goose-install goose-add goose-up goose-status goose-down
.PHONY: squawk-install squawk

bin-deps: .vendor.protogen
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(BIN_DIR) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${PROTO_GENERATE_PATH}
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go --go_out=${PROTO_GENERATE_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc --go-grpc_out=${PROTO_GENERATE_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway --grpc-gateway_out ${PROTO_GENERATE_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2 --openapiv2_out=${PROTO_GENERATE_PATH} \
		--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${PROTO_GENERATE_PATH}" \
		./api/manager-service/v1/manager-service.proto

.vendor.protogen: .vendor.protogen/google/protobuf .vendor.protogen/google/api .vendor.protogen/protoc-gen-openapiv2/options .vendor.protogen/validate

.vendor.protogen/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor.protogen/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor.protogen/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor.protogen/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp

clean:
	rm -rf $(BIN_DIR) godepgraph.png