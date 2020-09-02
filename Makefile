VERSION=1.0.7-dev
BUILDTIME=$(shell date +"%Y-%m-%dT%T%z")
LDFLAGS= -ldflags '-X github.com/bold-commerce/protoc-gen-struct-transformer/generator.version=$(VERSION) -X github.com/bold-commerce/protoc-gen-struct-transformer/generator.buildTime=$(BUILDTIME)'

.PHONY: re-generate-example generate install build version setup

re-generate-example: re-generate-example-v1 re-generate-example-v2

re-generate-example-v1:
	protoc \
		--proto_path=$(GOPATH)/pkg/mod/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=false,helper-package=helpers,goimports=true:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/message.proto

re-generate-example-v2:
	protoc \
		--proto_path=$(GOPATH)/pkg/mod/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=false,helper-package=helpers,goimports=true:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/v2/message.proto

generate: version re-generate-example

re-generate-example-debug: re-generate-example-debug-v1 re-generate-example-debug-v2

re-generate-example-debug-v1:
	protoc \
		--proto_path=$(GOPATH)/pkg/mod/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=true,helper-package=helpers,goimports=true:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/message.proto

re-generate-example-debug-v2:
	protoc \
		--proto_path=$(GOPATH)/pkg/mod/github.com/gogo:. \
		--struct-transformer_out=package=transform,debug=true,helper-package=helpers,goimports=true:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./example/v2/message.proto

generate-debug: version re-generate-example-debug

generate-annotations:
	protoc \
		--proto_path=$(GOPATH)/pkg/mod/github.com/gogo:. \
		--gogofaster_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options:. \
		./options/annotations.proto

install: setup
	go install $(LDFLAGS)

build: OUTPUT=.
build: setup
	go build $(LDFLAGS) -o $(OUTPUT)

version:
	protoc-gen-struct-transformer --version

setup:
	go mod download
	go mod verify
