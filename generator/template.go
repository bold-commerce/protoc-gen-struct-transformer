package generator

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"text/template/parse"
)

func at(t *template.Template) (string, *parse.Tree) {
	return t.Name(), t.Tree
}

func mt(name, tpl string, inc ...*template.Template) *template.Template {
	t := template.New(name).Funcs(funcMap)

	for _, v := range inc {
		if _, err := t.AddParseTree(at(v)); err != nil {
			log.Fatalln("unreachable")
		}
	}
	return template.Must(t.Parse(tpl))
}

var (
	funcMap = template.FuncMap{
		"formatField":          formatField,
		"formatOneofInitField": formatOneofInitField,
	}

	funcNameT = mt("FuncName", `{{- .SrcFn }}To{{ .DstFn }}`)
	srcParamT = mt("SrcParam", `{{- if .SrcPref }}{{- .SrcPref }}.{{ end }}{{ .Src }}, opts ...TransformParam`)
	dstParamT = mt("DstParam", `{{- if .DstPref }}{{- .DstPref }}.{{ end }}{{ .Dst }}`)
	ptrValT   = mt("PtrValName", `{{- if .Swapped -}} ValPtr {{- else -}} PtrVal {{- end }}`)
	ptrT      = mt("ptr", `{{ if .Ptr -}} Ptr {{- else -}} Val {{- end }}`)
	ptrOnlyT  = mt("ptrOnly", `{{ if .Ptr -}} Ptr {{- end }}`)
	starT     = mt("star", `{{ if .Ptr -}} * {{- end }}`)

	ptr2ptrT = mt("ptr2ptr", `func {{ template "FuncName" . }}Ptr(src *{{ template "SrcParam" . }}) *{{ template "DstParam" . }} {
	if src == nil {
		return nil
	}

	d := {{ template "FuncName" . }}(*src, opts...)
	return &d
}`, funcNameT, srcParamT, dstParamT)

	ptr2valT = mt("ptr2val", `func {{ template "FuncName" . }}PtrVal(src *{{ template "SrcParam" . }}) {{ template "DstParam" . }} {
	if src == nil {
		return {{ template "DstParam" . }}{}
	}

	return {{ template "FuncName" . }}(*src, opts...)
}`, funcNameT, srcParamT, dstParamT)

	val2ptrT = mt("val2ptr", `func {{ template "FuncName" . }}ValPtr(src {{ template "SrcParam" . }}) *{{ template "DstParam" . }} {
	d := {{ template "FuncName" . }}(src, opts...)
	return &d
}`, funcNameT, srcParamT, dstParamT)

	val2valT = mt("val2val", `func {{ template "FuncName" . }}(src {{ template "SrcParam" . }}) {{ template "DstParam" . }} {
	s := {{ template "DstParam" . }}{
		{{- with $R := . }}
			{{- range $f := .Fields}}
			{{ formatField $f $R.Swapped $R.DstPref }}
			{{- end -}}
		{{- end }}
	}

	applyOptions(opts...)

{{- with $R := . }}
{{ range $f := .Fields }}
{{ formatOneofInitField $f $R.Swapped }}
{{- end -}}
{{- end }}
	return s
}`, funcNameT, srcParamT, dstParamT)

	lst2lstT = mt("lst2lst", `func {{ template "FuncName" . }}{{ template "ptr" . }}List(src []{{ template "star" . }}{{ template "SrcParam" . }}) []{{ template "star" . }}{{ template "DstParam" . }} {
	resp := make([]{{ template "star" . }}{{ template "DstParam" . }}, len(src))

	for i, s := range src {
		resp[i] = {{ template "FuncName" . }}{{ template "ptrOnly" . }}(s, opts...)
	}

	return resp
}`, funcNameT, ptrT, srcParamT, starT, dstParamT, ptrOnlyT)

	ptrlst2ptrlstT = mt("ptrlst2ptrlst", `{{ template "lst2lst" .P true }}`, lst2lstT, funcNameT, ptrT, starT, srcParamT, dstParamT, ptrOnlyT)

	vallst2vallstT = mt("vallst2vallst", `{{ template "lst2lst" . }}`, lst2lstT, funcNameT, ptrT, starT, srcParamT, dstParamT, ptrOnlyT)

	ptrlst2vallstT = mt("ptrlst2vallst", `func {{ template "FuncName" . }}{{ template "PtrValName" . }}List(src []{{ .SrcPointer }}{{ template "SrcParam" . }}) []{{ .DstPointer }}{{ template "DstParam" . }} {
	resp := make([]{{ .DstPointer }}{{ template "DstParam" . }}, len(src))

	for i, s := range src {
		{{- if .DstPointer  }}
		g := {{ template "FuncName" . }}(s, opts...)
		resp[i] = &g
		{{ else }}
		resp[i] = {{ template "FuncName" . }}(*s)
		{{ end -}}
	}

	return resp
}`, funcNameT, ptrValT, srcParamT, dstParamT)

	ptr2vallstT = mt("ptr2vallst", `// {{ template "FuncName" . }}List is DEPRECATED. Use {{ template "FuncName" . }}{{ template "PtrValName" . }}List instead.
func {{ template "FuncName" . }}List(src []{{ .SrcPointer }}{{ template "SrcParam" . }}) []{{ .DstPointer }}{{ template "DstParam" . }} {
	return {{ template "FuncName" . }}{{ template "PtrValName" . }}List(src)
}`, funcNameT, ptrValT, srcParamT, dstParamT)

	tpls = []*template.Template{
		funcNameT, srcParamT, dstParamT, ptrValT, ptrT, ptrOnlyT, starT, ptr2ptrT,
		ptr2valT, val2ptrT, val2valT, lst2lstT, ptrlst2ptrlstT, vallst2vallstT,
		ptrlst2vallstT, ptr2vallstT,
	}

	// Executed with Data struct.
	oneFuncitonSetT = `{{- template "ptr2ptr" . }}

{{ template "ptrlst2ptrlst" . }}

{{ template "ptr2val" . }}

{{ template "ptrlst2vallst" . }}

{{ template "ptr2vallst" . }}

{{ template "val2val" . }}

{{ template "val2ptr" . }}

{{ template "vallst2vallst" . }}

`

	oneofT = `
type Oneof{{ .Decl }} interface {
	GetStringValue() string
	GetInt64Value() int64
}

func {{ .ProtoType }}To{{ .GoType }}(src Oneof{{ .Decl }}) string {
	if s := src.GetStringValue(); s != "" {
		return s
	}

	if i := src.GetInt64Value(); i != 0 {
		return strconv.FormatInt(i, 10)
	}

	return "<nil>"
}

func {{ .GoType }}To{{ .ProtoType }}(s string, dst *{{ .ProtoPackage }}.{{ .ProtoType }}, v string) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil  || v == "v2"{
		dst.{{ .Decl }} = &{{ .ProtoPackage }}.{{ .ProtoType }}_StringValue{StringValue: s}
		return
	}

	dst.{{ .Decl }} = &{{ .ProtoPackage }}.{{ .ProtoType }}_Int64Value{Int64Value: i}
	return
}

`

	optionsT = `var version string

// TransformParam is a function option type.
type TransformParam func()

// WithVersion sets global version variable.
func WithVersion(v string) TransformParam {
	return func() {
		version = v
	}
}

func applyOptions(opts ...TransformParam) {
	for _, o := range opts {
		o()
	}
}

`
)

// templateWithHelpers initializes main oneFuncitonSetT template with given
// name, adds there sub-templates and maps functions into template.
func templateWithHelpers(name string) (*template.Template, error) {
	t := template.
		New(name).
		Funcs(funcMap)

	for _, v := range tpls {
		if _, err := t.AddParseTree(at(v)); err != nil {
			return nil, err
		}
	}

	return t.Parse(oneFuncitonSetT)
}

// Field represents one structure field.
type Field struct {
	// Field name in Go structure.
	Name string
	// Field name in .proto file.
	ProtoName string
	// Field type in .proto file.
	ProtoType string
	// Name of function which is used for converting proto field into Go one.
	ProtoToGoType string
	// Name of function which is used for converting Go field into proto one.
	GoToProtoType string
	// True if field in model is a pointer.
	GoIsPointer bool
	// True if field in .proto file has an option gogoproto.nullable = false
	ProtoIsPointer bool
	// It true, field GoToProtoType and ProtoToGoType functions will be used
	// with prefix.
	UsePackage bool
	OneofDecl  string
	Opts       string
}

// IsOneof returns true if Field has non-empty OneOf declaration.
func (f Field) IsOneof() bool {
	return f.OneofDecl != ""
}

// name based on swapped flag return Name or ProtoName for current Field.
func (f Field) name(swapped bool) string {
	if swapped {
		return f.Name
	}
	return f.ProtoName
}

// convertFunc based on swapped flag and value of Field properties returns a type name
// for current Field.
func (f Field) convertFunc(swapped bool) string {
	out := f.ProtoToGoType
	if swapped {
		out = f.GoToProtoType
	}

	if f.GoIsPointer && f.ProtoIsPointer {
		out += "Ptr"
	}
	if !f.GoIsPointer && !f.ProtoIsPointer {
		out += ""
	}
	list := strings.HasSuffix(out, "List")
	if list {
		out = strings.TrimSuffix(out, "List")
	}

	suffix := ""

	if !f.GoIsPointer && f.ProtoIsPointer {
		if swapped {
			suffix = "ValPtr"
		} else {
			suffix = "PtrVal"
		}
	}

	if f.GoIsPointer && !f.ProtoIsPointer {
		if swapped {
			suffix = "PtrVal"
		} else {
			suffix = "ValPtr"
		}
	}

	if suffix != "" {
		out += suffix
		if list {
			out += "List"
		}
	}

	return out
}

// formatOneofField returns text representation of Oneof field in structure for
// template.
//
// This function is mapped into template. See funcMap variable for details.
func formatOneofField(f Field, swapped bool, pref string) string {
	if !f.IsOneof() {
		return fmt.Sprintf("/* field %q is not Oneof field*/", f.Name)
	}

	out := "src." + f.ProtoName
	if swapped {
		out = fmt.Sprintf("&%s.%s{}", pref, f.ProtoType)
	} else {
		if f.ProtoToGoType != "" {
			out = fmt.Sprintf("%s(src.%s)", f.ProtoToGoType, f.ProtoName)
		}
	}

	return out
}

// formatOneofInitField return text representation for filling up initialized
// oneof fields.
//
// This function is mapped into template. See funcMap variable for details.
func formatOneofInitField(f Field, swapped bool) string {
	if !swapped || !f.IsOneof() {
		return ""
	}
	return fmt.Sprintf(" %s(src.%s, s.%s, version)", f.GoToProtoType, f.Name, f.ProtoName)

}

func formatComplexField(f Field, swapped bool) string {
	if f.ProtoToGoType != "" {
		return fmt.Sprintf(" %s(src.%s %s)", f.convertFunc(swapped), f.name(swapped), f.Opts)
	}

	return fmt.Sprintf("src.%s", f.name(swapped))
}

// formatField returns a string with appropriate field convert functions for
// using in template.
func formatField(f Field, swapped bool, pref string) string {
	left := f.name(!swapped)

	right := ""
	if f.IsOneof() {
		right = formatOneofField(f, swapped, pref)
	} else {
		right = formatComplexField(f, swapped)
	}

	return fmt.Sprintf("%s: %s,", left, right)
}

// OneofData contains info about OneOf fields.
//
//	message TheOne{  <= OneofType
//	  oneof strint { <= OneofDecl
//	    string string_value = 1;
//	    int64 int64_value = 2;
//	  }
//	}
type OneofData struct {
	// Package name form proto file.
	ProtoPackage string
	// Custom data type from proto file which is used for define oneof field.
	ProtoType string
	// Go destination type.
	GoType    string
	Decl      string
	OneofDecl string
}

// Data contains data for fill out template.
type Data struct {
	// Prefix for source structure.
	SrcPref string
	// Source structure name.
	Src string
	// Left (source) part of transform function name.
	SrcFn string
	// Contains "*" if source structure is a pointer.
	SrcPointer string
	// Prefix for destination structure.
	DstPref string
	// Destination structure name.
	Dst string
	// Right (destination) part of transform function name.
	DstFn string
	// Contains "*" if destination structure is a pointer.
	DstPointer string
	// Field list of structure.
	Fields []Field
	// If true Fields.GoToProtoType will be used instead of Fileds.ProtoToGoType.
	Swapped bool
	// Is not empty, package name will be used as prefix for helper functions,
	// such as TimeToNullTime, StringToStirnPtr etc.
	HelperPackage string
	// Ptr is used in template for indication of pointer usage.
	Ptr bool
}

// swap swaps source and destination parameters for using in reverse functions.
func (d *Data) swap() {
	d.SrcPref, d.DstPref = d.DstPref, d.SrcPref
	d.Src, d.Dst = d.Dst, d.Src
	d.SrcFn, d.DstFn = d.DstFn, d.SrcFn
	d.SrcPointer, d.DstPointer = d.DstPointer, d.SrcPointer
	d.Swapped = !d.Swapped
}

// P sets Ptr flag of Data structure. Used inside template. Should be exported
// in template case.
func (d Data) P(t bool) Data {
	d.Ptr = t
	return d
}
