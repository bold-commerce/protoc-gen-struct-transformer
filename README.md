# Transformation function generator for gRPC.

[![Build Status](https://travis-ci.com/bold-commerce/protoc-gen-struct-transformer.svg?branch=master)](https://travis-ci.com/bold-commerce/protoc-gen-struct-transformer)
[![GoDoc](https://godoc.org/github.com/bold-commerce/protoc-gen-struct-transformer?status.svg)](https://godoc.org/github.com/bold-commerce/protoc-gen-struct-transformer)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bold-commerce/protoc-gen-struct-transformer?sort=semver)](https://github.com/bold-commerce/protoc-gen-struct-transformer/releases)
[![BSD-3-Clause](https://img.shields.io/github/license/bold-commerce/protoc-gen-struct-transformer)](./LICENSE)

<!-- vim-markdown-toc GFM -->

* [Quick presentation](#quick-presentation)
* [Overview](#overview)
* [How to use](#how-to-use)
  * [Installation](#installation)
    * [Homebrew](#homebrew)
    * [go get](#go-get)
  * [Add options to *.proto file](#add-options-to-proto-file)
  * [Run protoc](#run-protoc)
  * [Use generated functions in your gRPC server implementation.](#use-generated-functions-in-your-grpc-server-implementation)
  * [CLI parameters](#cli-parameters)
* [Troubleshooting](#troubleshooting)
  * [make generate returns an error](#make-generate-returns-an-error)
    * ["protobuf@v1.3.1/gogoproto/gogo.proto" was not found or had errors.](#protobufv131gogoprotogogoproto-was-not-found-or-had-errors)

<!-- vim-markdown-toc -->

## Quick presentation
[Speakerdeck](https://speakerdeck.com/ekhabarov/protoc-gen-struct-transformer)

## Overview
[Protocol buffers complier](https://github.com/protocolbuffers/protobuf) `protoc` generated structures based on message
definition in `*.proto` file. It's possible to use these generated structures
directly, but it's better to have clear separation between transport level
(gRPC) and business logic with its own structures. In this case you have to
convert protobuf structures into business logic structures and vice versa.

`protoc-gen-struct-transformer` is a plugin for `protoc` which generates functions
for structure transformation.

Let's look at simple example.

Source proto file:
```proto
// message.proto
syntax = "proto3";
package messages;

message Product {
  int32 id = 1;
  string name = 2;
}
```

Command `protoc --gogofaster_out=. message.proto` will generate `message.pb.go` with
following structure:
```go
type Product struct {
  Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
  Name string `protobuf:"bytes,2,opt,name=name,json=name,proto3" json:"name,omitempty"`
}
```
and let's suppose you service has as `repo` package with `ProductModel` struct inside:

```go
type ProductModel struct {
  ID   int    `db:"id" json:"id"`
  Name string `db:"name" json:"name"`
}
```
In order to publish data from the repo to API you have to convert `ProductModel`
to `Product`, for saving data which hit API you have to make back transformation.

```go
func ProductModelToProduct(m repo.ProductModel) proto.Product {
 return proto.Product {
    Id:   m.ID,
    Name: m.Name,
  }
}

func ProductToProductModel(p proto.Product) repo.ProductModel {
  return repo.ProductModel {
      ID:   p.Id,
      Name: p.Name,
  }
}
```
List of function type should be generated:

Source   | Destination | Suffix name
---------|-------------|---
*proto   | *model      | `Ptr`
[]*proto | []*model    | `PtrList`
*proto   | model       | `PtrVal`
[]*proto | []model     | `PtrValList`
proto    | model       |
proto    | \*model     | `ValPtr`
[]proto  | []model     | `ValList`
*model   | *proto      | `Ptr`
[]*model | []*proto    | `PtrList`
*model   | proto       | `PtrVal`
[]*model | []proto     | `PtrValList`
model    | proto       |
model    | \*proto     | `ValPtr`
[]model  | []proto     | `ValList`

function name has a format `<Source>To<Destination><Suffix>`.

For instance, function which converts list of pointers to Product into list of
Product values will be named `PbToProductPtrValList`,
where
* `Pb` is a replacement for proto message
* `Products` is a name of model structure
* `PtrValList` is a suffix pointed that convertion is made from slice of pointer to slice of values.

Full set of function for Product message will be as:
```go
func PbToProductPtr(src *example.Product, opts ...TransformParam) *model.Product
func PbToProductPtrList(src []*example.Product, opts ...TransformParam) []*model.Product
func PbToProductPtrVal(src *example.Product, opts ...TransformParam) model.Product
func PbToProductPtrValList(src []*example.Product, opts ...TransformParam) []model.Product
func PbToProductList(src []*example.Product, opts ...TransformParam) []model.Product
func PbToProduct(src example.Product, opts ...TransformParam) model.Product
func PbToProductValPtr(src example.Product, opts ...TransformParam) *model.Product
func PbToProductValList(src []example.Product, opts ...TransformParam) []model.Product
func ProductToPbPtr(src *model.Product, opts ...TransformParam) *example.Product
func ProductToPbPtrList(src []*model.Product, opts ...TransformParam) []*example.Product
func ProductToPbPtrVal(src *model.Product, opts ...TransformParam) example.Product
func ProductToPbValPtrList(src []model.Product, opts ...TransformParam) []*example.Product
func ProductToPbList(src []model.Product, opts ...TransformParam) []*example.Product
func ProductToPb(src model.Product, opts ...TransformParam) example.Product
func ProductToPbValPtr(src model.Product, opts ...TransformParam) *example.Product
func ProductToPbValList(src []model.Product, opts ...TransformParam) []example.Product
```

where
* `example` is a package generated by `protoc-gen-go` or `protoc-gen-gogo` plugin
* `model` is a package which contains manually created models structures.

Full example you can find in [example](./example) directory.

## How to use

### Installation
I assume you already have `protoc` installed.

First of all, it's necessary to install plugin itself, it's just a binary file,
which should be placed into $PATH to be available for `protoc`.

#### Homebrew

```shell
% brew tap bold-commerce/tap
% brew install protoc-gen-struct-transformer
```

#### go get
If you're going to make changes to plugin, use `go get ...` or `git clone ...`

```shell
% export GO111MODULE=on
% go get -u -d github.com/bold-commerce/protoc-gen-struct-transformer
% cd $GOPATH/src/github.com/bold-commerce/protoc-gen-struct-transformer
// make changes
% go install
```

Next, we need `protoc-gen-go` plugin (or `protoc-gen-gogofaster` if you use
`gogo` specific options) which creates `*.pb.go` file.
```shell
go get -u github.com/golang/protobuf/protoc-gen-go
// or
go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
```

### Add options to *.proto file
To configure plugin you have to use **file level** options listed below. The
plugin will not process file without these options.
```proto
// This import allows to use this options.
// Relatively to import_path: github.com/bold-commerce/protoc-gen-struct-transformer.
// See Makefile for mapping details.
import "options/annotations.proto";

// Go package name which contains business logic structures.
option (transformer.go_repo_package) = "models";
// Go package name with protobuf generated srtuctures. Could be equal to
// options go_package.
option (transformer.go_protobuf_package) = "example";
// Path to source file with Go structures which will be used as destination.
option (transformer.go_models_file_path) = "example/model/model.go";
```
as well as **message level** option
```proto
// Name of structure from business logic package. This option links business
// logic and generated structure.
message Product {
  option (transformer.go_struct) = "ProductModel";
  // ...
}
```
options above are minimal requirement for use this plugin.

Also plugin has additional **field level** options:

```proto
message Product {
  // SomeField will not be added to transformation function.
  string some_field = 4 [ (transformer.skip) = true ];
  // "map_as" option is used in cases when protoc-gen-go* plugin creates
  // "unpredictable" field name, i.e. by default protoc-gen-go* converts
  // protobuf named writen in snake_case into CamelCase, but for fields like
  // map_field_1 this rule has aa exception, in pb.go file it will be
  // "MapField_1" instead of "MapField1".
  // "map_to" options is used when you need to map current message field to
  // field in model with arbitrary name.
  // Both options "map_as" and "map_to" can be used independently.
  string map_field_1 = 6 [ (transformer.map_as) = "MapField_1", (transformer.map_to) = "MapField1"];
  // "custom" allows to use custom transformers for fields, which require extended transformation
  // The plugin won't generate methods for this field,
  // but rather expect it to be in the same package with the transformer file
  CustomType custom_field [(transformer.custom) = true]
}
```
### Run protoc
```shell
protoc \
  --proto_path=github.com/gogo:. \
  --go_out=Moptions/annotations.proto=github.com/bold-commerce/protoc-gen-struct-transformer/options,plugins=grpc:. \
  --struct-transformer_out=package=transform:. \
  ./message.proto
```
this command generates two files:
* `message.pb.go` contains auto-generated structures.
* `transform/message_transformer.go` contains transformation functions.

by default `message_transformer.go` does not contain imports. To add imports
run `protoc` with:
```shell
  --struct-transformer_out=package=transform,goimports=true:. \
```

### Use generated functions in your gRPC server implementation.
```go
func (s *server) CreateProduct(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	p, err := s.svc.Create(ctx, transform.PbToProduct(req.Product))
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		Product: transform.ProductToPb(p),
	}, nil
}
```

### CLI parameters
```
Usage of protoc-gen-struct-transformer:
  -debug
        Add debug information to generated file.
  -goimports
        Perform goimports on generated file.
  -helper-package string
        Package name for helper functions.
  -package string
        Package name for generated functions. (default "fallback")
  -use-package-in-path
        If true, package parameter will be used in path for output file. (default true)
  -version
        Print current version.
```
## Troubleshooting

### make generate returns an error
#### "protobuf@v1.3.1/gogoproto/gogo.proto" was not found or had errors.
`gogo.proto` file which is used for gogo-specific options is imported from
go modules cache. In order to fill out the cache run:

```shell
% export GO111MODULE=on
% go build
```
and run `make generate` again.
