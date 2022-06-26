ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), see https://golangci-lint.run/usage/install/#local-installation")
endif

.PHONY: format lint test build install freeze freeze-upgrade

default: all

all: install lint test build

format:
	$(info ******************** checking formatting ********************)
	go fmt ./...

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test:
	$(info ******************** running tests ********************)
	go test -cover -v ./...

build:
	$(info ******************** building bin ********************)
	goreleaser build --snapshot --rm-dist --single-target
	find dist -name seki

install:
	$(info ******************** downloading dependencies ********************)
	go mod download

freeze:
	$(info ******************** freeze dependencies ********************)
	go mod tidy && go mod verify

freeze-upgrade:
	$(info ******************** upgrade dependencies ********************)
	go get -u ./... && go mod tidy && go mod verify