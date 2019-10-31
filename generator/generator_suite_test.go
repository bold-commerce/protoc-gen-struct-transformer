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
		"StringField":  source.FieldInfo{Type: "string"},
		"BoolField":    source.FieldInfo{Type: "bool"},
		"IntField":     source.FieldInfo{Type: "int"},
		"Int32Field":   source.FieldInfo{Type: "int32"},
		"Int64Field":   source.FieldInfo{Type: "int64"},
		"UintField":    source.FieldInfo{Type: "uint"},
		"Uint32Field":  source.FieldInfo{Type: "uint32"},
		"Uint64Field":  source.FieldInfo{Type: "uint64"},
		"Float32Field": source.FieldInfo{Type: "float32"},
		"Float64Field": source.FieldInfo{Type: "float64"},
		"TimeField":    source.FieldInfo{Type: "time.Time"},
		"TimePtrField": source.FieldInfo{Type: "*time.Time"},
		"PkgTypeField": source.FieldInfo{Type: "pkg.Type"},
		"ProtoField":   source.FieldInfo{Type: "proto.FieldType"},

		"StringFieldPtr":  source.FieldInfo{Type: "string", IsPointer: true},
		"BoolFieldPtr":    source.FieldInfo{Type: "bool", IsPointer: true},
		"IntFieldPtr":     source.FieldInfo{Type: "int", IsPointer: true},
		"Int32FieldPtr":   source.FieldInfo{Type: "int32", IsPointer: true},
		"Int64FieldPtr":   source.FieldInfo{Type: "int64", IsPointer: true},
		"UintFieldPtr":    source.FieldInfo{Type: "uint", IsPointer: true},
		"Uint32FieldPtr":  source.FieldInfo{Type: "uint32", IsPointer: true},
		"Uint64FieldPtr":  source.FieldInfo{Type: "uint64", IsPointer: true},
		"Float32FieldPtr": source.FieldInfo{Type: "float32", IsPointer: true},
		"Float64FieldPtr": source.FieldInfo{Type: "float64", IsPointer: true},
		"TimeFieldPtr":    source.FieldInfo{Type: "time.Time", IsPointer: true},
		"TimePtrFieldPtr": source.FieldInfo{Type: "*time.Time", IsPointer: true},
		"PkgTypeFieldPtr": source.FieldInfo{Type: "pkg.Type", IsPointer: true},
		"ProtoFieldPtr":   source.FieldInfo{Type: "proto.FieldType", IsPointer: true},
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
