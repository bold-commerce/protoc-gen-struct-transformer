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
)

var _ = Describe("Message", func() {
	var (
		messagesData = source.StructureList{
			"msg1": goStruct,
		}
	)

	Describe("processMessage", func() {

		DescribeTable("check result",
			func(msg *descriptor.DescriptorProto, dstStruct string, expFields []Field, expSructName string, expError error) {
				if msg != nil && dstStruct != "" {
					err := proto.SetExtension(msg.Options, options.E_GoStruct, sp(dstStruct))
					Expect(err).NotTo(HaveOccurred())
				}

				fields, structName, err := processMessage(nil, msg, subm, messagesData, false)
				if expError == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(MatchError(expError))
				}

				Expect(fields).To(Equal(expFields))
				Expect(structName).To(Equal(expSructName))

			},
			Entry("Nil message", nil, "", nil, "", newLoggableError("message is nil")),

			Entry("Message without fields", &descriptor.DescriptorProto{
				Name:    sp("Msg1"),
				Field:   nil,
				Options: &descriptor.MessageOptions{},
			}, "msg1", []Field{}, "msg1", nil),

			Entry("Message with non_existent field", &descriptor.DescriptorProto{
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
			}, "msg1", nil, "", errors.New(`field "NotExists" not found in destination structure`)),

			Entry("Message with fields", &descriptor.DescriptorProto{
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
			}, "msg1", []Field{
				{
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
				},
			}, "msg1", nil),
		)
	})

})
