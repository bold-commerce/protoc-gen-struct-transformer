// NOTE: This file is NOT autogenerated and contains custom transformers,
//       which are used in the message_transformer.go

package transform

import (
	"strconv"

	exampleV2 "github.com/bold-commerce/protoc-gen-struct-transformer/example/v2"
)

// PbCustomTypeToStringPtrVal is an example of the custom transformer from Pb to go
func PbCustomTypeToStringPtrVal(src *exampleV2.CustomType, opts ...TransformParam) string {
	applyOptions(opts...)

	if version == "v2" {
		return src.Value
	}

	return ""
}

// StringToPbCustomTypeValPtr is an example of the custom transformer from go to Pb
func StringToPbCustomTypeValPtr(src string, opts ...TransformParam) *exampleV2.CustomType {
	applyOptions(opts...)

	if version == "v2" {
		return &exampleV2.CustomType{
			Value: src,
		}
	}

	return nil
}

// PbTheOneToStringPtrVal is a custom transformer for the OneOf object
func PbTheOneToStringPtrVal(src *exampleV2.TheOne, opts ...TransformParam) string {
	if s := src.GetStringValue(); s != "" {
		return s
	}

	if i := src.GetInt64Value(); i != 0 {
		return strconv.FormatInt(i, 10)
	}

	return "<nil>"
}

// StringToPbTheOneValPtr is the custom transformer for the OneOf object with versions
func StringToPbTheOneValPtr(s string, opts ...TransformParam) *exampleV2.TheOne {
	applyOptions(opts...)

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil || version == "v2" {
		return &exampleV2.TheOne{TheDecl: &exampleV2.TheOne_StringValue{StringValue: s}}
	}

	return &exampleV2.TheOne{TheDecl: &exampleV2.TheOne_Int64Value{Int64Value: i}}
}