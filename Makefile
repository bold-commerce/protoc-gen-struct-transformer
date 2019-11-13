VERSION=1.0.1-dev
BUILDTIME=$(shell date +"%Y-%m-%dT%T%z")
LDFLAGS= -ldflags '-X github.com/bold-commerce/protoc-gen-struct-transformer/generator.version=$(VERSION) -X github.com/bold-commerce/protoc-gen-struct-transformer/generator.buildTime=$(BUILDTIME)'

.PHONY: re-generate-example
re-generate-example:
	protoc \
		--proto_path=vendor/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=false:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/message.proto

.PHONY: imports
imports:
	$(GOBIN)/goimports -w example/transform/message_transformer.go

.PHONY: generate
generate: re-generate-example imports

install:
	go install $(LDFLAGS)

build: OUTPUT=.
build:
	go build $(LDFLAGS) -o $(OUTPUT)
