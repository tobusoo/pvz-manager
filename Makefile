include .env

PROTO_GENERATE_PATH=$(CURDIR)/pkg

GOCYCLO_PATH=$(shell go env GOPATH)/bin/gocyclo
GOCOGNIT_PATH=$(shell go env GOPATH)/bin/gocognit
GODEPGRAPH_PATH=$(shell go env GOPATH)/bin/godepgraph
GOOSEE_PATH=$(shell go env GOPATH)/bin/goose

MIGRATIONS_PATH=./migrations

THRESHOLD=5

BIN_DIR=$(CURDIR)/bin
SERVICE_NAME=manager

SERVICE_PATH_SRC=./cmd/$(SERVICE_NAME)_service
CLI_PATH_SRC=./cmd/$(SERVICE_NAME)_cli
NOTIFIER_PATH_SRC=./cmd/notifier

SERVICE_PATH_BIN=$(BIN_DIR)/$(SERVICE_NAME)_service
CLI_PATH_BIN=$(BIN_DIR)/$(SERVICE_NAME)_cli
NOTIFIER_BIN=$(BIN_DIR)/notifier

SERVICE_DOCKERFILE_PATH=build/dev/$(SERVICE_NAME)_service/Dockerfile
SERVICE_DOCKER_CONTAINER_NAME=manager-service-image:1.0.0
NOTIFIER_DOCKERFILE_PATH=build/dev/notifier/Dockerfile
NOTIFIER_DOCKER_CONTAINER_NAME=notifier-image:1.0.0
DOCKER_DEV_COMPOSE_PATH=build/dev/docker-compose.yml
DOCKER_TEST_COMPOSE_PATH=build/test/docker-compose.yml

.PHONY: all mkdir-bin run build tidy clean gocyclo gocognit test coverage
.PHONY: unit-test integration-test integration-test-db e2e-test benchmark

all: bin-deps generate build compose-up goose-up

run-cli: build
	$(CLI_PATH_BIN)

unit-test:
	@echo "Unit Tests:"
	@go test ./internal/usecase/ -coverprofile=coverage_usecase.out
	@go test ./internal/app/manager_service/ -coverprofile=coverage_manager.out

integration-test:
	docker-compose -f $(DOCKER_TEST_COMPOSE_PATH) up -d
	@echo "Sleeping 4 seconds for postgreSQL preparation"
	@sleep 4
	@$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_TEST_DSN) up
	@POSTGRESQL_TEST_DSN=${POSTGRESQL_TEST_DSN} go test -v -coverpkg=./internal/storage/postgres \
		-coverprofile=coverage_storage_postgres.out \
		./tests/integration/storage_db/integration_test.go
	@go test -v ./tests/integration/kafka/integration_test.go
	docker-compose -f $(DOCKER_TEST_COMPOSE_PATH) down

integration-test-intergnal:
	@echo "Integration Tests:"
	@go test -coverpkg=./internal/storage/storage_json -coverprofile=coverage_storage.out \
		./tests/integration/storage_json/integration_test.go

e2e-test: build
	@echo "E2E Tests:"
	@go test tests/e2e_test.go

test: unit-test integration-test-intergnal integration-test
	@echo "mode: set" > coverage.out
	@tail -n +2 coverage_usecase.out >> coverage.out
	@tail -n +2 coverage_storage.out >> coverage.out
	@tail -n +2 coverage_storage_postgres.out >> coverage.out
	@tail -n +2 coverage_manager.out >> coverage.out
	@rm coverage_usecase.out coverage_storage.out coverage_storage_postgres.out coverage_manager.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html 

benchmark:
	@go test -bench=. -benchtime=10x  -benchmem ./benchmark/storage_test.go

build: dependancy-install gocyclo gocognit mkdir-bin $(SERVICE_PATH_BIN) $(CLI_PATH_BIN) $(NOTIFIER_BIN)

mkdir-bin:
	@mkdir -p bin

$(SERVICE_PATH_BIN): mkdir-bin
	go build -o $(SERVICE_PATH_BIN) $(SERVICE_PATH_SRC)

$(CLI_PATH_BIN): mkdir-bin
	go build -o $(CLI_PATH_BIN) $(CLI_PATH_SRC)

$(NOTIFIER_BIN): mkdir-bin
	go build -o $(NOTIFIER_BIN) $(NOTIFIER_PATH_SRC)

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

docker-build: docker-build-service docker-build-notifier

docker-build-service:
	docker build --no-cache -f $(SERVICE_DOCKERFILE_PATH) . -t $(SERVICE_DOCKER_CONTAINER_NAME)

docker-build-notifier:
	docker build --no-cache -f $(NOTIFIER_DOCKERFILE_PATH) . -t $(NOTIFIER_DOCKER_CONTAINER_NAME)

compose-up:
	docker compose -f $(DOCKER_DEV_COMPOSE_PATH) up --detach

compose-down:
	docker compose -f $(DOCKER_DEV_COMPOSE_PATH) down

compose-stop:
	docker compose -f $(DOCKER_DEV_COMPOSE_PATH) stop

compose-start:
	docker compose -f $(DOCKER_DEV_COMPOSE_PATH) start

compose-ps:
	docker compose -f $(DOCKER_DEV_COMPOSE_PATH) ps

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN_LOCAL) create rename_me sql

goose-up:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN_LOCAL) up

goose-down:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN_LOCAL) down

goose-status:
	$(GOOSEE_PATH) -dir $(MIGRATIONS_PATH) postgres $(POSTGRESQL_DSN_LOCAL) status

squawk-install:
	npm install -g squawk-cli

squawk:
	squawk ./migrations/* --exclude=ban-drop-table

.PHONY: depgraph compose-up compose-down compose-stop compose-start goose-install goose-add goose-up goose-status goose-down
.PHONY: squawk-install squawk

$(BIN_DIR)/protoc-gen-go:
	@GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(BIN_DIR)/protoc-gen-go-grpc:
	@GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$(BIN_DIR)/protoc-gen-grpc-gateway:
	@GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

$(BIN_DIR)/protoc-gen-openapiv2:
	@GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

$(BIN_DIR)/protoc-gen-validate:
	@GOBIN=$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@latest

$(BIN_DIR)/statik:
	@GOBIN=$(BIN_DIR) go install github.com/rakyll/statik@latest

bin-deps: vendor.protogen $(BIN_DIR)/protoc-gen-go $(BIN_DIR)/protoc-gen-go-grpc $(BIN_DIR)/protoc-gen-grpc-gateway $(BIN_DIR)/protoc-gen-grpc-gateway $(BIN_DIR)/protoc-gen-openapiv2 $(BIN_DIR)/protoc-gen-validate $(BIN_DIR)/statik

generate:
	@mkdir -p ${PROTO_GENERATE_PATH}
	@protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go --go_out=${PROTO_GENERATE_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc --go-grpc_out=${PROTO_GENERATE_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway --grpc-gateway_out ${PROTO_GENERATE_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2 --openapiv2_out=${PROTO_GENERATE_PATH} \
		--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${PROTO_GENERATE_PATH}" \
		./api/manager-service/v1/manager-service.proto

vendor.protogen: vendor.protogen/google/protobuf vendor.protogen/google/api vendor.protogen/protoc-gen-openapiv2/options vendor.protogen/validate

vendor.protogen/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

vendor.protogen/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

vendor.protogen/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor.protogen/google
	mv vendor.protogen/googleapis/google/api vendor.protogen/google
	rm -rf vendor.protogen/googleapis

vendor.protogen/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p vendor.protogen/validate
	mv vendor.protogen/tmp/validate vendor.protogen/
	rm -rf vendor.protogen/tmp

clean:
	rm -rf $(BIN_DIR) godepgraph.png coverage*.out