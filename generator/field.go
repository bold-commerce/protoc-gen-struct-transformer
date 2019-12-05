package generator

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bold-commerce/protoc-gen-struct-transformer/options"
	"github.com/bold-commerce/protoc-gen-struct-transformer/source"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/iancoleman/strcase"
	pkgerrors "github.com/pkg/errors"
)

// lastName splits string by "." and returns last part.
func lastName(s string) string {
	splt := strings.Split(s, ".")
	return splt[len(splt)-1]
}

// wktgoogleProtobufTimestamp returns *Field created out of
// google.protobuf.Timestamp protobuf field.
func wktgoogleProtobufTimestamp(pname, gname string, gf source.FieldInfo, pnullable bool) *Field {
	p2g := ""
	g2p := ""

	if gf.Type != "time.Time" {
		g := strcase.ToCamel(strings.Replace(gf.Type, ".", "", -1))
		p := "Time"

		if pnullable {
			p += "Ptr"
		}

		if gf.IsPointer {
			g += "Ptr"
		}

		p2g = fmt.Sprintf("%sTo%s", p, g)
		g2p = fmt.Sprintf("%sTo%s", g, p)
	}

	return &Field{
		Name:          gname,
		ProtoName:     pname,
		ProtoToGoType: p2g,
		GoToProtoType: g2p,
		UsePackage:    p2g != "",
	}
}

// wktgoogleProtobufString returns *Field created out of
// google.protobuf.StringValue field.
func wktgoogleProtobufString(pname, gname, ftype string) *Field {
	g := strcase.ToCamel(strings.Replace(ftype, ".", "", -1))
	p := "StringValue"

	return &Field{
		Name:          gname,
		ProtoName:     pname,
		ProtoToGoType: fmt.Sprintf("%sTo%s", p, g),
		GoToProtoType: fmt.Sprintf("%sTo%s", g, p),
		UsePackage:    true,
	}
}

// processSubMessage processes sub messages of current message. Sub message is
// a message type which is used as field type.
//
// In the next example message B is a current message and message A is sub
// message.
//
// message A {}
// message B { A a_field = 1; }
func processSubMessage(w io.Writer,
	fdp *descriptor.FieldDescriptorProto,
	pname, gname, pbtype string,
	mo MessageOption,
	goStructFields source.Structure,
) (*Field, error) {

	if fdp == nil {
		return nil, errors.New("input field is nil")
	}

	if fdp.Name == nil {
		return nil, errors.New("input field name is nil")
	}

	tpl := "%sTo%s"
	pb := "Pb"

	p2g := ""
	g2p := ""

	if mo != nil {
		if mo.OneofDecl() != "" {
			pb = strcase.ToCamel(goStructFields[gname].Type)
		} else {
			pb, pbtype = mo.Target(), pb
		}

		pb = strcase.ToCamel(pb)
	}

	if l := fdp.Label; l != nil && *l == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		tpl += "List"
		if g, ok := goStructFields[gname]; ok {
			pb = strcase.ToCamel(g.Type)
		}
	}

	// embedded fields
	fname := gname
	if isEmbed := extractEmbedOption(fdp.Options); isEmbed {
		// if sub message is embedded use type name as field name.
		fname = pb
		pb = strcase.ToCamel(pb)
	}

	if ln := lastName(pbtype); strings.Contains(pbtype, ".") {
		pbtype = strcase.ToCamel(ln)
	}
	isNullable := extractNullOption(fdp)

	p2g = fmt.Sprintf(tpl, pbtype, pb)
	g2p = fmt.Sprintf(tpl, pb, pbtype)

	f := &Field{
		Name:           strcase.ToCamel(fname),
		ProtoName:      strcase.ToCamel(*fdp.Name),
		ProtoType:      pbtype,
		ProtoToGoType:  p2g,
		GoToProtoType:  g2p,
		Opts:           ", opts...",
		ProtoIsPointer: isNullable,
	}

	if fm, ok := goStructFields[gname]; ok {
		if mo == nil {
			return nil, errors.New("mo is nil")
		}
		f.GoIsPointer = fm.IsPointer
		f.OneofDecl = mo.OneofDecl()
	}

	return f, nil
}

// processSimpleField processes fields of basic types such as int, string and
// so on.
func processSimpleField(w io.Writer, pname, gname string, ftype *descriptor.FieldDescriptorProto_Type, sf source.FieldInfo) (*Field, error) {

	sf.Type = strcase.ToCamel(strings.Replace(sf.Type, ".", "", -1)) // pkg.Type => PkgType
	t := types[*ftype]

	// sf: NullsString, pbType: , goType: string, ft: TYPE_STRING, name: Tags, pbaname: Tags
	p(w, "// sf: %#v, pbType: %q, goType: %q, ft: %q, pname: %q, gname: %q\n",
		sf, t.pbType, t.goType, ftype, pname, gname)

	sft := strings.ToLower(sf.Type)
	tpb := strings.ToLower(t.pbType)
	tgo := strings.ToLower(t.goType)

	f := &Field{
		Name:      gname,
		ProtoName: pname,
	}

	switch true {

	case (sft == tpb && tpb != "") || (sft == tgo && tpb == ""): // equal types
		f.ProtoToGoType = ""
		f.GoToProtoType = ""

	case sft != tgo:
		p := t.pbType
		if p == "" {
			p = t.goType
		}

		f.ProtoToGoType = fmt.Sprintf("%sTo%s", strcase.ToCamel(p), sf.Type)
		f.GoToProtoType = fmt.Sprintf("%sTo%s", sf.Type, strcase.ToCamel(p))
		f.UsePackage = true

	case sft != tpb:
		p(w, "// sft: %s, tpb: %s\n", sft, tpb)
		f.ProtoToGoType = sft
		f.GoToProtoType = tpb

	default:
		f.ProtoToGoType = t.pbType
		f.GoToProtoType = t.goType
		f.UsePackage = t.usePackage
	}

	return f, nil
}

// processField returns filled Field struct for template.
func processField(
	w io.Writer,
	fdp *descriptor.FieldDescriptorProto,
	subMessages MessageOptionList,
	goStructFields source.Structure,
) (*Field, error) {
	// If field has transformer.skip == true, it will be not processed.
	if skip := extractSkipOption(fdp.Options); skip {
		return nil, newLoggableError("field skipped: %s", *fdp.Name)
	}

	mapTo, err := getStringOption(fdp.Options, options.E_MapTo)
	if _, ok := err.(errOptionNotExists); err != nil && err != ErrNilOptions && !ok {
		return nil, pkgerrors.Wrap(err, "mapTo option")
	}

	// if field has an options map_as then overwrite fieldName which is pbname
	mapAs, err := getStringOption(fdp.Options, options.E_MapAs)
	if _, ok := err.(errOptionNotExists); err != nil && err != ErrNilOptions && !ok {
		return nil, pkgerrors.Wrap(err, "mapAs option")
	}

	pname, gname := prepareFieldNames(*fdp.Name, mapAs, mapTo)

	// check if field exists in destination/Go structure.
	gf, ok := goStructFields[gname]
	if !ok {
		// do not check for embedded fields.
		if isEmbed := extractEmbedOption(fdp.Options); !isEmbed {
			return nil, pkgerrors.Wrap(errors.New("field not found in destination structure"), gname)
		}
	}

	p(w, "\n\n// ===============================\n")
	if oi := fdp.OneofIndex; oi != nil {
		p(w, "// fdp.OneofIndex: %#v\n\n", *oi)
	}

	if tn := fdp.TypeName; tn != nil {
		p(w, "// fdp.TypeName: %#v\n\n", *tn)
	}
	if opt := fdp.Options; opt != nil {
		p(w, "// fdp.Options: %s\n\n", strings.Replace(fmt.Sprintf("%#v", opt), "\n", "", -1))
	}
	p(w, "// fdp.Name: %q, mapAs: %q, mapTo: %q\n", *fdp.Name, mapAs, mapTo)

	// Process subMessages. For details see comments for the TypeName.
	if typ := fdp.TypeName; *fdp.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE && typ != nil {
		t := *typ
		switch t {
		case ".google.protobuf.Timestamp":
			isNullable := extractNullOption(fdp)
			return wktgoogleProtobufTimestamp(pname, gname, gf, isNullable), nil
		case ".google.protobuf.StringValue":
			return wktgoogleProtobufString(pname, gname, gf.Type), nil
		}

		// Submessage has a name like ".package.type", 1: removes first ".".
		mo, _ := subMessages[t[1:]]
		// TODO(ekhabarov): pass gf instead of goStructFields
		return processSubMessage(w, fdp, pname, gname, t, mo, goStructFields)
	}

	return processSimpleField(w, pname, gname, fdp.Type, gf)
}

// abbreviationUpper checks a incoming string for equality and suffixes, if it
// exists it will be converted to uppercase.
// For instance, identifier fields in models often have a name like SomeID, with
// capitalized "ID", while protobuf auto-generated structures use names like
// "SomeId".
// TODO(ekhabarov): Add cli parameter for such mapping.
func abbreviationUpper(name string) string {
	abbreviation := []string{"Id", "Sku", "Url"}

	for _, a := range abbreviation {
		if name == a {
			return strings.ToUpper(a)
		}

		if strings.HasSuffix(name, a) {
			return strings.TrimSuffix(name, a) + strings.ToUpper(a)
		}
	}

	return name
}

// prepareFieldNames returns names Protobuf  and Go for field, considering
// map_to/map_as options and abbreviation rules.
func prepareFieldNames(fname, mapAs, mapTo string) (string, string) {
	pname := strcase.ToCamel(fname)
	if mapAs != "" {
		pname = mapAs
	}

	gname := abbreviationUpper(pname)
	if mapTo != "" {
		gname = mapTo
	}

	return pname, gname
}
