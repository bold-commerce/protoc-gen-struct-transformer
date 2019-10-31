package generator

import "fmt"

// MessageOption represents protobuf message options.
type MessageOption interface {
	// Returns model name from go_struct option.
	Target() string
	// Full returns FQTN of proto type, such as "google.protobuf.Timestamp" and
	// so on.
	Full() string
	// If true, proto message has not target struct, i.e. go_struct option is
	// empty or not found.
	Omitted() bool
	// Returns Oneof message name.
	OneofDecl() string
}

// MessageOptionList is a list of proto message option. Map key is a message
// name with Full Qualified Type Name (FQTN), format like
// ".google.protobuf.Timestamp".
type MessageOptionList map[string]MessageOption

func (sol MessageOptionList) String() string {
	s := "\n"
	for k, v := range sol {
		s += fmt.Sprintf("// %q: target: %q, Omitted: %t, OneofDecl: %q\n",
			k, v.Target(), v.Omitted(), v.OneofDecl())
	}

	return s
}

type messageOption struct {
	// Value of transformer.go_struct option.
	targetName string
	// Full message name e.g. svc.example.Address.
	fullName string
	// OneOf name.
	oneofDecl string
}

func (so messageOption) Target() string {
	return so.targetName
}

func (so messageOption) Full() string {
	return so.fullName
}

func (so messageOption) Omitted() bool {
	return so.Target() == ""
}

func (so messageOption) OneofDecl() string {
	return so.oneofDecl
}
