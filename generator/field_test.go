package generator

import (
	"errors"

	"github.com/bold-commerce/protoc-gen-struct-transformer/options"
	"github.com/bold-commerce/protoc-gen-struct-transformer/source"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	pkgerrors "github.com/pkg/errors"
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

					func(pname, gname, typ string, gp, pnullable bool, expected Field) {
						got := wktgoogleProtobufTimestamp(pname, gname, source.FieldInfo{Type: typ, IsPointer: gp}, pnullable)

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

					Entry("Field not found", "protoName", "name", "AnyGoType", false, false, Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "TimeToAnyGoType",
						GoToProtoType: "AnyGoTypeToTime",
						UsePackage:    true,
					}),
					Entry("String", "protoName", "name", "string", false, false, Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "TimeToString",
						GoToProtoType: "StringToTime",
						UsePackage:    true,
					}),
					Entry("Time", "protoName", "name", "time.Time", false, false, Field{
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

					Entry("Field not found", "protoName", "name", "AnyGoType", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "StringValueToAnyGoType",
						GoToProtoType: "AnyGoTypeToStringValue",
						UsePackage:    true,
					}),
					Entry("String", "protoName", "name", "int64", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "StringValueToInt64",
						GoToProtoType: "Int64ToStringValue",
						UsePackage:    true,
					}),
					Entry("pkg.Type", "protoName", "name", "pkg.Type", Field{
						Name:          "name",
						ProtoName:     "protoName",
						ProtoToGoType: "StringValueToPkgType",
						GoToProtoType: "PkgTypeToStringValue",
						UsePackage:    true,
					}),
				)
			})
		})
	})

	Describe("ProcessSubMessages", func() {

		var (
			protoField         = "proto_field"
			protoFieldTypeName = "CustomType"
			goField            = "StringField"
			labelRepeated      = descriptor.FieldDescriptorProto_LABEL_REPEATED
		)

		DescribeTable("check result",
			func(fdp *descriptor.FieldDescriptorProto, pname, gname, pbType string, mo MessageOption, custom bool, expected *Field) {
				got, err := processSubMessage(nil, fdp, pname, gname, pbType, mo, goStruct, custom)
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

			Entry("Int64", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", mo, false, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "Pb",
				ProtoToGoType:  "PbToMoTarget",
				GoToProtoType:  "MoTargetToPb",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}),

			Entry("Custom field", &descriptor.FieldDescriptorProto{Name: &protoField, TypeName: &protoFieldTypeName}, protoField, goField, "int64", mo, true, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "PbCustomType",
				ProtoToGoType:  "PbCustomTypeToString",
				GoToProtoType:  "StringToPbCustomType",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}),

			Entry("With messageOption and empty oneof", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", mo, false, &Field{
				Name:           "StringField",
				ProtoName:      "ProtoField",
				ProtoType:      "Pb",
				ProtoToGoType:  "PbToMoTarget",
				GoToProtoType:  "MoTargetToPb",
				GoIsPointer:    false,
				ProtoIsPointer: true,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           ", opts...",
			}),

			Entry("With messageOption and non-empty oneof", &descriptor.FieldDescriptorProto{Name: &protoField}, protoField, goField, "int64", moWithOneOf, false, &Field{
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
				protoField, goField, "full.type", moWithOneOf, false,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "Type",
					ProtoToGoType:  "TypeToString",
					GoToProtoType:  "StringToType",
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
				protoField, goField, "string", mo, false,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "Pb",
					ProtoToGoType:  "PbToStringList",
					GoToProtoType:  "StringToPbList",
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
				protoField, goField, "string", mo, false,
				&Field{
					Name:           "StringField",
					ProtoName:      "ProtoField",
					ProtoType:      "Pb",
					ProtoToGoType:  "PbToStringList",
					GoToProtoType:  "StringToPbList",
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
			pint32  = descriptor.FieldDescriptorProto_TYPE_INT32
			pint64  = descriptor.FieldDescriptorProto_TYPE_INT64
			pstring = descriptor.FieldDescriptorProto_TYPE_STRING
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

			Entry("Same types: int64", "Abc", "Abc", &pint64, goStruct["Int64Field"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "",
					GoToProtoType:  "",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Same types: string", "Abc", "Abc", &pstring, goStruct["StringField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "",
					GoToProtoType:  "",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Similar type: int64 <=> int", "Abc", "Abc", &pint64, goStruct["IntField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "int",
					GoToProtoType:  "int64",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     false,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Different types: int32", "Abc", "Abc", &pint32, goStruct["StringField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "Int32ToString",
					GoToProtoType:  "StringToInt32",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     true,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Different types: int64", "Abc", "Abc", &pint64, goStruct["StringField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "Int64ToString",
					GoToProtoType:  "StringToInt64",
					GoIsPointer:    false,
					ProtoIsPointer: false,
					UsePackage:     true,
					OneofDecl:      "",
					Opts:           "",
				}),

			Entry("Custom Go type into basic proto type", "Abc", "Abc", &pstring, goStruct["PkgTypeField"],
				&Field{
					Name:           "Abc",
					ProtoName:      "Abc",
					ProtoType:      "",
					ProtoToGoType:  "StringToPkgType",
					GoToProtoType:  "PkgTypeToString",
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
			Entry("Field: ID", "ID", "", "", "ID", "ID"),
			Entry("Field: id", "id", "", "", "Id", "ID"),
			Entry("Field: Id", "Id", "", "", "Id", "ID"),
			Entry("Field: iD", "iD", "", "", "ID", "ID"),
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
					// Check just a message here.
					Expect(err.Error()).To(Equal(expectedErr.Error()))
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
				GoToProtoType:  "",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           "",
			}, nil),

			Entry("int64: capitalized ID", &descriptor.FieldDescriptorProto{
				Name:     sp("ID"),
				TypeName: sp("int64"),
				Type:     &typInt64,
				Options:  &descriptor.FieldOptions{},
			}, false, false, &Field{
				Name:           "ID",
				ProtoName:      "ID",
				ProtoType:      "",
				ProtoToGoType:  "",
				GoToProtoType:  "",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     false,
				OneofDecl:      "",
				Opts:           "",
			}, nil),

			Entry("int64: id", &descriptor.FieldDescriptorProto{
				Name:     sp("id"),
				TypeName: sp("int64"),
				Type:     &typInt64,
				Options:  &descriptor.FieldOptions{},
			}, false, false, &Field{
				Name:           "ID",
				ProtoName:      "Id",
				ProtoType:      "",
				ProtoToGoType:  "",
				GoToProtoType:  "",
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
			}, false, false, nil, pkgerrors.Wrap(errors.New("field not found in destination structure"), "NotExists")),

			Entry("embed", &descriptor.FieldDescriptorProto{
				Name:     sp("PkgTypeField"),
				TypeName: sp(".PkgType"),
				Type:     &typMessage,
				Options:  &descriptor.FieldOptions{},
			}, false, true, &Field{
				Name:           "PkgField",
				ProtoName:      "PkgTypeField",
				ProtoType:      "Pb",
				ProtoToGoType:  "PbToPkgField",
				GoToProtoType:  "PkgFieldToPb",
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
				ProtoToGoType:  "StringValueToString",
				GoToProtoType:  "StringToStringValue",
				GoIsPointer:    false,
				ProtoIsPointer: false,
				UsePackage:     true,
				OneofDecl:      "",
				Opts:           "",
			}, nil),
		)

	})

})
