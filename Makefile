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

install: VERSION=1.0.0-dev
install: SHA1=$(shell git rev-parse --short HEAD)
install: BUILDTIME=$(shell date +"%Y-%m-%d_%H:%M:%S")
install:
	go install -ldflags '-X github.com/bold-commerce/protoc-gen-struct-transformer/generator.gitHash=$(SHA1) -X github.com/bold-commerce/protoc-gen-struct-transformer/generator.version=$(VERSION) -X github.com/bold-commerce/protoc-gen-struct-transformer/generator.buildTime=$(BUILDTIME)'
