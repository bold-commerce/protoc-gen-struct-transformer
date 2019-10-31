package generator

import (
	"errors"

	"github.com/bold-commerce/protoc-gen-struct-transformer/options"
	"github.com/bold-commerce/protoc-gen-struct-transformer/source"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Field", func() {

	Describe("abbreviation", func() {

		Context("when an abbreviation words passed", func() {

			DescribeTable("returns an abbreviation in uppercase",
				func(abbr, result string) {
					got := abbreviationUpper(abbr)
					Expect(got).To(Equal(result))
				},

				Entry("ID", "Id", "ID"),
				Entry("URL", "Url", "URL"),
				Entry("SKU", "Sku", "SKU"),

				Entry("NewID", "NewId", "NewID"),
				Entry("ProdURL", "ProdUrl", "ProdURL"),
				Entry("SKU", "SomeSku", "SomeSKU"),

				Entry("NonAbbr", "NonAbbr", "NonAbbr"),
			)

		})
	})

	Describe("Well-known types", func() {

		Describe("google.protobuf.Timestamp", func() {

			Context("when fields if of type", func() {

				DescribeTable("check Field stuct",

					func(pname, gname, ftype string, expected Field) {
						got := wktgoogleProtobufTimestamp(pname, gname, ftype)

						Expect(*got).To(MatchAllFields(Fields{
							"Name":           Equal(expected.Name),
							"ProtoName":      Equal(expected.ProtoName),
							"ProtoToGoType":  Equal(expected.ProtoToGoType),
							"GoToProtoType":  Equal(expected.GoToProtoType),
							"ProtoType":      Equal(expected.ProtoType),
							"GoIsPointer":    Equal(expected.GoIsPointer),
							"ProtoIsPointer": Equal(expected.ProtoIsPointer),
							"UsePackage":     Equal(expected.UsePackage),
							"OneofDecl":      Equal(expected.OneofDecl),
							"Opts":           Equal(expected.Opts),
						}))
					},

					Entry("Field not found", "protoName", "name", "AnyType", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "AnyTypeToTimePtr",
						GoToProtoType: "TimePtrToAnyType",
						UsePackage:    true,
					}),
					Entry("String", "protoName", "name", "string", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "StringToTimePtr",
						GoToProtoType: "TimePtrToString",
						UsePackage:    true,
					}),
					Entry("Time", "protoName", "name", "time.Time", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "",
						GoToProtoType: "",
						UsePackage:    false,
					}),
				)
			})
		})

		Describe("google.Protobuf.StringValue", func() {

			Context("when fields if of type", func() {

				DescribeTable("check Field stuct",

					func(pname, gname, ftype string, expected Field) {
						got := wktgoogleProtobufString(pname, gname, ftype)

						Expect(*got).To(MatchAllFields(Fields{
							"Name":           Equal(expected.Name),
							"ProtoName":      Equal(expected.ProtoName),
							"ProtoToGoType":  Equal(expected.ProtoToGoType),
							"GoToProtoType":  Equal(expected.GoToProtoType),
							"ProtoType":      Equal(expected.ProtoType),
							"GoIsPointer":    Equal(expected.GoIsPointer),
							"ProtoIsPointer": Equal(expected.ProtoIsPointer),
							"UsePackage":     Equal(expected.UsePackage),
							"OneofDecl":      Equal(expected.OneofDecl),
							"Opts":           Equal(expected.Opts),
						}))
					},

					Entry("Field not found", "protoName", "name", "AnyType", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "AnyTypeToStringValue",
						GoToProtoType: "StringValueToAnyType",
						UsePackage:    true,
					}),
					Entry("String", "protoName", "name", "int64", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "Int64ToStringValue",
						GoToProtoType: "StringValueToInt64",
						UsePackage:    true,
					}),
					Entry("Time", "protoName", "name", "pkg.Type", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "PkgTypeToStringValue",
						GoToProtoType: "StringValueToPkgType",
						UsePackage:    true,
					}),
				)
			})
		})
	})

	Describe("ProcessSubMessages", func() {

		var (
			protoField    = "proto_field"
			goField       = "StringField"
			labelRepeated = descriptor.FieldDescriptorProto_LABEL_REPEATED
		)

		DescribeTable("check result",
			func(fdp *descriptor.FieldDescriptorProto, pname, gname, pbType string, mo MessageOption, expected *Field) {
				got, err := processSubMessage(nil, fdp, pname, gname, pbType, mo, goStruct)
				Expect(err).NotTo(HaveOccurred())

				Expect(*got).To(MatchAllFields(Fields{
					"Name":           Equal(expected.Name),
					"ProtoName":      Equal(expected.ProtoName),
					"ProtoToGoType":  Equal(expected.ProtoToGoType),
					"GoToProtoType":  Equal(expected.GoToProtoType),
					"ProtoType":      Equal(expected.ProtoType),
					"GoIsPointer":    Equal(expected.GoIsPointer),
					"ProtoIsPointer": Equal(expected.ProtoIsPointer),
					"UsePackage":     Equal(expected.UsePackage),
					"OneofDecl":      Equal(expected.OneofDecl),
					"Opts":           Equal(expected.Opts),
				}))
			},

			Entry("Int64", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", mo, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "moTarget",
				ProtoToGoType:  "moTargetToPb",
				GoToProtoType:  "PbTomoTarget",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}),

			Entry("With messageOption and empty oneof", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", mo, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "moTarget",
				ProtoToGoType:  "moTargetToPb",
				GoToProtoType:  "PbTomoTarget",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}),

			Entry("With messageOption and non-empty oneof", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", moWithOneOf, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "int64",
				ProtoToGoType:  "int64ToString",
				GoToProtoType:  "StringToint64",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "oneofField",
				Opts:           ", opts...",
			}),

			Entry("With messageOption, empty oneof, and fqdn type name",
				&descriptor.FieldDescriptorProto{
					Name: &protoField,
				},
				protoField, goField, "full.type", moWithOneOf,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "type",
					ProtoToGoType:  "typeToString",
					GoToProtoType:  "StringTotype", // TODO(ekhabarov): should be fixed. ToType
					GoIsPointer:    false,
					ProtoIsPointer: true,
					UsePackage:     false,
					OneofDecl:      "oneofField",
					Opts:           ", opts...",
				}),

			Entry("Repeated field",
				&descriptor.FieldDescriptorProto{
					Name:  &protoField,
					Label: &labelRepeated,
				},
				protoField, goField, "string", mo,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "FieldType",
					ProtoToGoType:  "FieldTypeToPbList",
					GoToProtoType:  "PbToFieldTypeList",
					GoIsPointer:    false,
					ProtoIsPointer: true,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           ", opts...",
				}),

			Entry("Repeated field when name field found in target struct.",
				&descriptor.FieldDescriptorProto{
					Name:  &protoField,
					Label: &labelRepeated,
				},
				protoField, goField, "string", mo,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "FieldType",
					ProtoToGoType:  "FieldTypeToPbList",
					GoToProtoType:  "PbToFieldTypeList",
					GoIsPointer:    false,
					ProtoIsPointer: true,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           ", opts...",
				}),
		)
	})

	Describe("ProcessSimpleField", func() {

		var (
			pint64 = descriptor.FieldDescriptorProto_TYPE_INT64
		)

		DescribeTable("check result",
			func(pname, gname string, ftype *descriptor.FieldDescriptorProto_Type, sf source.FieldInfo, expected *Field) {
				got, err := processSimpleField(nil, pname, gname, ftype, sf)
				Expect(err).NotTo(HaveOccurred())

				Expect(*got).To(MatchAllFields(Fields{
					"Name":           Equal(expected.Name),
					"ProtoName":      Equal(expected.ProtoName),
					"ProtoToGoType":  Equal(expected.ProtoToGoType),
					"GoToProtoType":  Equal(expected.GoToProtoType),
					"ProtoType":      Equal(expected.ProtoType),
					"GoIsPointer":    Equal(expected.GoIsPointer),
					"ProtoIsPointer": Equal(expected.ProtoIsPointer),
					"UsePackage":     Equal(expected.UsePackage),
					"OneofDecl":      Equal(expected.OneofDecl),
					"Opts":           Equal(expected.Opts),
				}))

			},

			Entry("Same types", "Abc", "Abc", &pint64, goStruct["Int64Field"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "",
					GoToProtoType:  "int64",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Different types", "Abc", "Abc", &pint64, goStruct["StringField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "StringToInt",
					GoToProtoType:  "IntToString",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     true,
					OneofDecl:      "",
					Opts:           "",
				}),
		)
	})

	Describe("prepareFieldNames", func() {

		DescribeTable("parameter combinations",
			func(fname, a, t, expectA, expectT string) {
				mapAs, mapTo := prepareFieldNames(fname, a, t)
				Expect(mapAs).To(Equal(expectA))
				Expect(mapTo).To(Equal(expectT))
			},

			Entry("Empty mapping", "proto_field_name", "", "", "ProtoFieldName", "ProtoFieldName"),
			Entry("MapAs without mapTo", "proto_field_name", "map_as", "", "map_as", "map_as"),
			Entry("MapTo without mapAs", "proto_field_name", "", "map_to", "ProtoFieldName", "map_to"),
			Entry("MapTo and mapAs", "proto_field_name", "map_as", "map_to", "map_as", "map_to"),
		)

	})

	Describe("processField", func() {

		DescribeTable("check result",
			func(f *descriptor.FieldDescriptorProto, skip, embed bool, expected *Field, expectedErr error) {

				err := proto.SetExtension(f.Options, options.E_Skip, bp(skip))
				Expect(err).NotTo(HaveOccurred())

				err = proto.SetExtension(f.Options, options.E_Embed, bp(embed))
				Expect(err).NotTo(HaveOccurred())

				field, err := processField(nil, f, subm, goStruct)
				if expectedErr == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(expectedErr))
				}

				if expectedErr == nil {
					Expect(*field).To(MatchAllFields(Fields{
						"Name":           Equal(expected.Name),
						"ProtoName":      Equal(expected.ProtoName),
						"ProtoToGoType":  Equal(expected.ProtoToGoType),
						"GoToProtoType":  Equal(expected.GoToProtoType),
						"ProtoType":      Equal(expected.ProtoType),
						"GoIsPointer":    Equal(expected.GoIsPointer),
						"ProtoIsPointer": Equal(expected.ProtoIsPointer),
						"UsePackage":     Equal(expected.UsePackage),
						"OneofDecl":      Equal(expected.OneofDecl),
						"Opts":           Equal(expected.Opts),
					}))
				}
			},

			Entry("int64", &descriptor.FieldDescriptorProto{
				Name:     sp("int64_field"),
				TypeName: sp("int64"),
				Type:     &typInt64,
				Options:  &descriptor.FieldOptions{},
			}, false, false, &Field{
				Name:           "Int64Field",
				ProtoName:      "Int64Field",
				ProtoType:      "",
				ProtoToGoType:  "",
				GoToProtoType:  "int64",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           "",
			}, nil),

			Entry("Skip", &descriptor.FieldDescriptorProto{
				Name:     sp("int64_field"),
				TypeName: sp("int64"),
				Type:     &typInt64,
				Options:  &descriptor.FieldOptions{},
			}, true, false, nil, newLoggableError("field skipped: int64_field")),

			Entry("Target field not found", &descriptor.FieldDescriptorProto{
				Name:     sp("not_exists"),
				TypeName: sp("int64"),
				Type:     &typInt64,
				Options:  &descriptor.FieldOptions{},
			}, false, false, nil, errors.New(`field "NotExists" not found in destination structure`)),

			Entry("embed", &descriptor.FieldDescriptorProto{
				Name:     sp("PkgTypeField"),
				TypeName: sp(".PkgType"),
				Type:     &typMessage,
				Options:  &descriptor.FieldOptions{},
			}, false, true, &Field{
				Name:           "PkgField",
				ProtoName:      "PkgTypeField",
				ProtoType:      "pkgField",
				ProtoToGoType:  "pkgFieldToPb",
				GoToProtoType:  "PbTopkgField",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}, nil),

			Entry("WKT: Timestamp", &descriptor.FieldDescriptorProto{
				Name:     sp("time_field"),
				TypeName: sp(".google.protobuf.Timestamp"),
				Type:     &typMessage,
				Options:  &descriptor.FieldOptions{},
			}, false, true, &Field{
				Name:           "TimeField",
				ProtoName:      "TimeField",
				ProtoType:      "",
				ProtoToGoType:  "",
				GoToProtoType:  "",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           "",
			}, nil),

			Entry("WKT: StringValue", &descriptor.FieldDescriptorProto{
				Name:     sp("string_field"),
				TypeName: sp(".google.protobuf.StringValue"),
				Type:     &typMessage,
				Options:  &descriptor.FieldOptions{},
			}, false, true, &Field{
				Name:           "StringField",
				ProtoName:      "StringField",
				ProtoType:      "",
				ProtoToGoType:  "StringToStringValue",
				GoToProtoType:  "StringValueToString",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     true,
				OneofDecl:      "",
				Opts:           "",
			}, nil),
		)

	})

})
