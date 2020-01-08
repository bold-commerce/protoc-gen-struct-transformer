# How to contribute to proto-gen-struct-transformer

Thank you for your contribution to this project. Next several steps describe
process of contribution:

- Please, open an issue first and describe what problem you are trying to solve.
- Make changes.
- Add test(s) for new code.
- If your changes modify plugin's output, please, add an appropriate example to `example` directory and re-generate it with `make generate`.
- Run `ginkgo -r -cover` on your feature branch and master branch. New feature should not decrease test coverage.
- Open PR on GitHub.

# Developers tools
- [Protocol buffers compiler (protoc)](https://github.com/protocolbuffers/protobuf) - Google's data interchange format.
- [protoc-gen-gogofaster](https://github.com/gogo/protobuf/tree/master/protoc-gen-gogofaster) - protoc plugin implements Go bindings for protocol buffers.
- [goimports](https://golang.org/x/tools/cmd/goimports) - Command goimports updates your Go import lines, adding missing ones and removing unreferenced ones.
- [Ginkgo](https://github.com/onsi/ginkgo#set-me-up) - BDD Testing Framework for Go.
