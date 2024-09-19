GOCYCLO_PATH=$(shell go env GOPATH)/bin/gocyclo
GOCOGNIT_PATH=$(shell go env GOPATH)/bin/gocognit
GODEPGRAPH_PATH=$(shell go env GOPATH)/bin/godepgraph

THRESHOLD=5

BIN_DIR=bin
APP_NAME=manager

APP_PATH_SRC=cmd/$(APP_NAME)/main.go
APP_PATH_BIN=$(BIN_DIR)/$(APP_NAME)

.PHONY: all mkdir-bin run build tidy clean gocyclo gocognit test coverage

all: build

run: build
	./$(APP_PATH_BIN)

test: build
	go test ./internal/usecase/ -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html 

build: dependancy-install mkdir-bin $(APP_PATH_BIN) gocyclo gocognit

mkdir-bin:
	@mkdir -p bin

$(APP_PATH_BIN): mkdir-bin
	go build -o $(APP_PATH_BIN) $(APP_PATH_SRC)

dependancy-update:
	@go get -u

dependancy-install:
	@go get ./...

tidy:
	@go mod tidy

gocyclo-install:
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gocognit-install:
	@go install github.com/uudashr/gocognit/cmd/gocognit@latest

gocyclo: gocyclo-install
	$(GOCYCLO_PATH) -over $(THRESHOLD) -ignore "_mock|_test" .

gocognit: gocognit-install
	$(GOCOGNIT_PATH) -over $(THRESHOLD) -ignore "_mock|_test" .

depgraph-install:
	@go install github.com/kisielk/godepgraph@latest

depgraph-build:
	$(GODEPGRAPH_PATH) -s $(APP_PATH_SRC) | dot -Tpng -o godepgraph.png

depgraph: depgraph-install depgraph-build

.PHONY: depgraph

clean:
	rm -rf $(BIN_DIR) godepgraph.png