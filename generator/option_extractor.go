package generator

import (
	"fmt"

	"github.com/bold-commerce/protoc-gen-struct-transformer/options"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// extractStructNameOption returns transformer.go_struct option value.
func extractStructNameOption(msg *descriptor.DescriptorProto) (string, error) {
	if msg == nil {
		return "", newLoggableError("message is nil")
	}

	if msg.Options == nil || !proto.HasExtension(msg.Options, options.E_GoStruct) {
		return "", newLoggableError("message %q has no option %q, skipped...", *msg.Name, options.E_GoStruct.Name)
	}

	ext, err := proto.GetExtension(msg.Options, options.E_GoStruct)
	if err != nil {
		return "", err
	}

	option, ok := ext.(*string)
	if !ok {
		return "", fmt.Errorf("extension is %T; want an *string", ext)
	}

	return *option, nil
}

// getStringOption return any option of string type for proto.Message. If
// option exists but has different type, function returns an error.
func getStringOption(m proto.Message, opt *proto.ExtensionDesc) (string, error) {
	if m == nil {
		return "", ErrNilOptions
	}

	if !proto.HasExtension(m, opt) {
		return "", newErrOptionNotExists(opt.Name)
	}

	ext, err := proto.GetExtension(m, opt)
	if err != nil {
		return "", err
	}

	option, ok := ext.(*string)
	if !ok {
		return "", fmt.Errorf("extension is %T; want an *string", ext)
	}

	return *option, nil
}

// getBoolOption return any option of bool type for proto.Message. If
// option exists but has different type, function returns false.
func getBoolOption(m proto.Message, opt *proto.ExtensionDesc) bool {
	if m == nil {
		return false
	}

	if !proto.HasExtension(m, opt) {
		return false
	}

	ext, err := proto.GetExtension(m, opt)
	if err != nil {
		return false
	}

	option, ok := ext.(*bool)
	if !ok {
		return false
	}

	return *option
}

// extractEmbedOption returns true if proto.Message has an option
// transformer.embed which equals to true.
func extractEmbedOption(m proto.Message) bool {
	return getBoolOption(m, options.E_Embed)
}

// extractSkipOption return value of transformer.skip option or false if
// option does not exist.
func extractSkipOption(m proto.Message) bool {
	return getBoolOption(m, options.E_Skip)
}

// extractNullOption returns true if Field has a gogoproto.nullable option which
// equals to true.
func extractNullOption(f *descriptor.FieldDescriptorProto) bool {
	return gogoproto.IsNullable(f)
}
