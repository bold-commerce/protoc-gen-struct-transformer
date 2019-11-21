VERSION=1.0.2-dev
BUILDTIME=$(shell date +"%Y-%m-%dT%T%z")
LDFLAGS= -ldflags '-X github.com/bold-commerce/protoc-gen-struct-transformer/generator.version=$(VERSION) -X github.com/bold-commerce/protoc-gen-struct-transformer/generator.buildTime=$(BUILDTIME)'

.PHONY: re-generate-example imports generate install build version

re-generate-example:
	protoc \
		--proto_path=vendor/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=false,helper-package=helpers:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/message.proto

imports:
	$(GOBIN)/goimports -w example/transform/message_transformer.go

generate: version re-generate-example imports

install:
	go install $(LDFLAGS)

build: OUTPUT=.
build:
	go build $(LDFLAGS) -o $(OUTPUT)

version:
	protoc-gen-struct-transformer --version
