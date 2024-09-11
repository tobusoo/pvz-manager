GOCYCLO_PATH=$(shell go env GOPATH)/bin/gocyclo
GOCOGNIT_PATH=$(shell go env GOPATH)/bin/gocognit
GODEPGRAPH_PATH=$(shell go env GOPATH)/bin/godepgraph

THRESHOLD=5

APP_NAME=manager

.PHONY: all run build tidy clean gocyclo gocognit

all: build

run: build
	./$(APP_NAME)

build: dependancy-install $(APP_NAME) gocyclo gocognit

$(APP_NAME):
	go build -o $(APP_NAME) main.go

dependancy-update:
	@go get -u

dependancy-install:
	@go get .

tidy:
	@go mod tidy

gocyclo-install:
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gocognit-install:
	@go install github.com/uudashr/gocognit/cmd/gocognit@latest

gocyclo: gocyclo-install
	$(GOCYCLO_PATH) -over $(THRESHOLD) .

gocognit: gocognit-install
	$(GOCOGNIT_PATH) -over $(THRESHOLD) .

depgraph-install:
	go install github.com/kisielk/godepgraph@latest

depgraph-build:
	$(GODEPGRAPH_PATH) -s main.go | dot -Tpng -o godepgraph.png

depgraph: depgraph-install depgraph-build

.PHONY: depgraph

clean:
	rm -rf $(APP_NAME) godepgraph.png