package generator

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

// processOneofFields adds function for first found field with OneOf
// declaration.
func processOneofFields(w io.Writer, data []*Data) error {
	added := map[string]struct{}{}

	for _, d := range data {
		if d == nil {
			continue
		}
		d.swap()

		for _, f := range d.Fields {
			if !f.IsOneof() {
				continue
			}

			pt := strings.Split(f.GoToProtoType, "To")[0]
			gt := strings.Split(f.ProtoToGoType, "To")[0]

			if _, ok := added[gt]; ok {
				continue
			}

			added[gt] = struct{}{}

			od := OneofData{
				ProtoType:    gt,
				ProtoPackage: d.SrcPref,
				GoType:       pt,
				Decl:         strcase.ToCamel(f.OneofDecl),
				OneofDecl:    "___decl___",
			}

			t, err := template.
				New("oneof" + f.ProtoName).
				Parse(tplOneof)
			if err != nil {
				return err
			}

			if err := t.Execute(w, od); err != nil {
				return err
			}
			// Add oneof function for first found field only.
			break
		}
	}
	return nil
}

// OptHelpers returns file content with optional functions for using options
// with transformations.
func OptHelpers(packageName string) string {
	w := output()
	fmt.Fprintln(w, "\npackage", packageName)
	fmt.Fprintln(w, tplOption)

	return w.String()
}
