package generator

import "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

type typeRel struct {
	pbType     string
	goType     string
	usePackage bool
}

// types contains protobuf types.
var types = map[descriptor.FieldDescriptorProto_Type]typeRel{
	descriptor.FieldDescriptorProto_TYPE_INT32:  typeRel{pbType: "int32", goType: "int"},
	descriptor.FieldDescriptorProto_TYPE_INT64:  typeRel{pbType: "int64", goType: "int"},
	descriptor.FieldDescriptorProto_TYPE_UINT32: typeRel{pbType: "uint32", goType: "uint"},
	descriptor.FieldDescriptorProto_TYPE_UINT64: typeRel{pbType: "uint64", goType: "uint"},
	descriptor.FieldDescriptorProto_TYPE_FLOAT:  typeRel{pbType: "", goType: "float32"},
	descriptor.FieldDescriptorProto_TYPE_DOUBLE: typeRel{pbType: "", goType: "float64"},
	descriptor.FieldDescriptorProto_TYPE_BOOL:   typeRel{pbType: "", goType: "bool"},
	descriptor.FieldDescriptorProto_TYPE_STRING: typeRel{pbType: "", goType: "string"},
}
