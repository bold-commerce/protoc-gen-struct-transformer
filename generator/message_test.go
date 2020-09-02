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
	pkgerrors "github.com/pkg/errors"
)

var _ = Describe("Message", func() {
	var (
		messagesData = source.StructureList{
			"msg1": goStruct,
		}
	)

	const (
		v1 = int32(1)
		v2 = int32(2)
	)

	Describe("processMessage", func() {

		DescribeTable("check result",
			func(msg *descriptor.DescriptorProto, dstStruct string, version int32, expFields []Field, expSructName string, expError error) {
				if msg != nil && dstStruct != "" {
					err := proto.SetExtension(msg.Options, options.E_GoStruct, sp(dstStruct))
					Expect(err).NotTo(HaveOccurred())
				}

				fConf := FileConfig{Version: version}

				fields, structName, err := processMessage(nil, msg, subm, messagesData, fConf)
				if expError == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					// Check just a message here.
					Expect(err.Error()).To(Equal(expError.Error()))
				}

				Expect(fields).To(Equal(expFields))
				Expect(structName).To(Equal(expSructName))

			},

			// Verison 1 tests

			Entry("Nil message - V1", nil, "", v1, nil, "", newLoggableError("message is nil")),

			Entry("Message without fields - V1", &descriptor.DescriptorProto{
				Name:    sp("Msg1"),
				Field:   nil,
				Options: &descriptor.MessageOptions{},
			}, "msg1", v1, []Field{}, "msg1", nil),

			Entry("Message with non_existent field - V1", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("not_exists"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v1, nil, "", pkgerrors.Wrap(errors.New("field not found in destination structure"), "NotExists")),

			Entry("Message with fields - V1", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("int64_field"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v1, []Field{
				{
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
				},
			}, "msg1", nil),

			Entry("Message with ID field - V1", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("ID"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v1, []Field{
				{
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
				},
			}, "msg1", nil),

			// Version 2 tests

			Entry("Nil message - V2", nil, "", v2, nil, "", newLoggableError("message is nil")),

			Entry("Message without fields", &descriptor.DescriptorProto{
				Name:    sp("Msg1"),
				Field:   nil,
				Options: &descriptor.MessageOptions{},
			}, "msg1", v2, []Field{}, "msg1", nil),

			Entry("Message with non_existent field - V2", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("not_exists"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v2, nil, "", pkgerrors.Wrap(errors.New("field not found in destination structure"), "NotExists")),

			Entry("Message with fields - V2", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("int64_field"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v2, []Field{
				{
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
				},
			}, "msg1", nil),

			Entry("Message with ID field - V2", &descriptor.DescriptorProto{
				Name: sp("Msg1"),
				Field: []*descriptor.FieldDescriptorProto{
					&descriptor.FieldDescriptorProto{
						Name:     sp("ID"),
						Number:   nil,
						Label:    nil,
						Type:     &typInt64,
						TypeName: nil, // sub message type
						Options:  &descriptor.FieldOptions{},
					},
				},
				Options: &descriptor.MessageOptions{},
			}, "msg1", v2, []Field{
				{
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
				},
			}, "msg1", nil),
		)
	})

})
