package generator

import (
	"io"

	"github.com/bold-commerce/protoc-gen-struct-transformer/source"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// processMessage processes each message regardless of contains it an options or
// it doesn't. It returns set of fields for template and destination structure
// name extracted from proto message go_struct option.
func processMessage(
	w io.Writer,
	msg *descriptor.DescriptorProto,
	subMessages map[string]MessageOption,
	str source.StructureList,
	fConf FileConfig,
) ([]Field, string, error) {

	structName, err := extractStructNameOption(msg)
	if err != nil {
		if msg != nil {
			for _, d := range msg.OneofDecl {
				p(w, "// Oneof: %#v\n\n", *d.Name)
			}
		}

		return nil, "", err
	}

	tsf, err := source.Lookup(str, structName)
	if err != nil {
		return nil, "", err
	}

	debugWriter := (io.Writer)(nil)
	if fConf.Debug {
		debugWriter = w
	}

	p(debugWriter, "%s", tsf)

	fields := []Field{}

	for _, f := range msg.Field {
		pf, err := processField(debugWriter, f, subMessages, tsf, fConf)
		if err != nil {
			if e, ok := err.(loggableError); ok {
				p(w, "// %s\n", e)
				continue
			}
			if err != ErrNilOptions {
				return nil, "", err
			}
			p(w, "// error: %s\n", err)
			continue
		}

		fields = append(fields, *pf)
	}

	return fields, structName, nil
}
