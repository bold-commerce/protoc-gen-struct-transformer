// Transformation function generator for gRPC.
//
// Overview
//
// Protocol buffers complier (protoc) https://github.com/protocolbuffers/protobuf
// generates structures based on message definition in *.proto file. It's
// possible to use these generated structures directly, but it's better to have
// clear separation between transport level (gRPC) and business logic with its
// own structures. In this case you have to convert generated structures into
// structures use in business logic and vice versa.
//
// See documentation and usage examples on https://github.com/bold-commerce/protoc-gen-struct-transformer/blob/master/README.md
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bold-commerce/protoc-gen-struct-transformer/generator"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"github.com/golang/protobuf/proto"
)

var (
	packageName       = flag.String("package", "fallback", "Package name for generated functions.")
	helperPackageName = flag.String("helper-package", "", "Package name for helper functions.")
	versionFlag       = flag.Bool("version", false, "Print current version.")
	debug             = flag.Bool("debug", false, "Add debug information to generated file.")
	usePackageInPath  = flag.Bool("use-package-in-path", true, "If true, package parameter will be used in path for output file.")
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(generator.Version())
		os.Exit(0)
	}

	var gogoreq plugin.CodeGeneratorRequest

	data, err := ioutil.ReadAll(os.Stdin)
	must(err)
	must(proto.Unmarshal(data, &gogoreq))

	// Convert incoming parameters into CLI flags.
	must(generator.SetParameters(flag.CommandLine, gogoreq.Parameter))

	resp := &plugin.CodeGeneratorResponse{}
	optPath := ""

	messages, err := generator.CollectAllMessages(gogoreq)
	must(err)

	for _, f := range gogoreq.ProtoFile {

		filename, content, err := generator.ProcessFile(f, packageName, helperPackageName, messages, *debug, *usePackageInPath)
		if err != nil {
			if err != generator.ErrFileSkipped {
				must(err)
			}
			continue
		}

		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(filename),
			Content: proto.String(content),
		})

		optPath = filename
	}

	if optPath != "" {
		optPath = filepath.Dir(optPath) + "/options.go"

		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(optPath),
			Content: proto.String(generator.OptHelpers(*packageName)),
		})
	}

	// Send back the results.
	data, err = proto.Marshal(resp)
	must(err)

	_, err = os.Stdout.Write(data)
	must(err)
}

func must(err error) {
	if err != nil {
		if *debug {
			log.Fatalf("%+v", err)
		} else {
			log.Fatalf("%v", err)
		}
	}
}
