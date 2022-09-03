ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), see https://golangci-lint.run/usage/install/#local-installation")
endif

ifeq (, $(shell which goreleaser))
$(warning "could not find goreleaser in $(PATH), see https://goreleaser.com/install/")
endif

ifeq (, $(shell which gotestsum))
$(warning "could not find gotestsum in $(PATH), see https://github.com/gotestyourself/gotestsum#install")
endif


.PHONY: format lint test build install freeze freeze-upgrade generate

default: all

all: install generate lint test build

format:
	$(info ******************** checking formatting ********************)
	go fmt ./...

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test:
	$(info ******************** running tests ********************)
	gotestsum

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

generate:
	$(info ******************** generating support files ********************)
	go generate ./...
