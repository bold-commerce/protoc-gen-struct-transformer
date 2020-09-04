package generator

import (
	"testing"

	"github.com/bold-commerce/protoc-gen-struct-transformer/source"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generator Suite")
}

var (
	typInt64   = descriptor.FieldDescriptorProto_TYPE_INT64
	typMessage = descriptor.FieldDescriptorProto_TYPE_MESSAGE

	sp = func(s string) *string {
		return &s
	}
	bp = func(b bool) *bool {
		return &b
	}

	// key - field name, value - field type
	// goStruct contains model structure fields.
	goStruct = map[string]source.FieldInfo{
		"ID":           {Type: "int64"},
		"StringField":  {Type: "string"},
		"BoolField":    {Type: "bool"},
		"IntField":     {Type: "int"},
		"Int32Field":   {Type: "int32"},
		"Int64Field":   {Type: "int64"},
		"UintField":    {Type: "uint"},
		"Uint32Field":  {Type: "uint32"},
		"Uint64Field":  {Type: "uint64"},
		"Float32Field": {Type: "float32"},
		"Float64Field": {Type: "float64"},
		"TimeField":    {Type: "time.Time"},
		"TimePtrField": {Type: "*time.Time"},
		"PkgTypeField": {Type: "pkg.Type"},
		"ProtoField":   {Type: "proto.FieldType"},

		"StringFieldPtr":  {Type: "string", IsPointer: true},
		"BoolFieldPtr":    {Type: "bool", IsPointer: true},
		"IntFieldPtr":     {Type: "int", IsPointer: true},
		"Int32FieldPtr":   {Type: "int32", IsPointer: true},
		"Int64FieldPtr":   {Type: "int64", IsPointer: true},
		"UintFieldPtr":    {Type: "uint", IsPointer: true},
		"Uint32FieldPtr":  {Type: "uint32", IsPointer: true},
		"Uint64FieldPtr":  {Type: "uint64", IsPointer: true},
		"Float32FieldPtr": {Type: "float32", IsPointer: true},
		"Float64FieldPtr": {Type: "float64", IsPointer: true},
		"TimeFieldPtr":    {Type: "time.Time", IsPointer: true},
		"TimePtrFieldPtr": {Type: "*time.Time", IsPointer: true},
		"PkgTypeFieldPtr": {Type: "pkg.Type", IsPointer: true},
		"ProtoFieldPtr":   {Type: "proto.FieldType", IsPointer: true},
	}

	mo = messageOption{
		targetName: "moTarget",
		fullName:   "full.name",
	}
	moWithOneOf = messageOption{
		targetName: "target",
		fullName:   "full.name",
		oneofDecl:  "oneofField",
	}
	moPkgField = messageOption{
		targetName: "pkgField",
		fullName:   "full.name",
	}

	subm = map[string]MessageOption{
		"FieldName": mo,
		"two":       moWithOneOf,
		"PkgType":   moPkgField,
	}
)
