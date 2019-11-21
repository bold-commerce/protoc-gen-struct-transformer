package source

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strconv"
)

// inspect is a function which is run for each node in source file. See go/ast
// package for details.
func inspect(output StructureList) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			// skip non-types
			return true
		}

		if spec.Type == nil {
			// skip empty types
			return true
		}

		s, ok := spec.Type.(*ast.StructType)
		if !ok {
			// skip non-struct types
			return true
		}

		structName := spec.Name.Name
		if _, ok := output[structName]; !ok {
			output[structName] = Structure{}
		}

		embeddedCounter := 0
		for _, field := range s.Fields.List {
			fname := "embedded_"
			// Embedded strcuts have no names.
			if field.Names != nil {
				fname = field.Names[0].Name
			} else {
				fname += strconv.Itoa(embeddedCounter)
				embeddedCounter++
			}

			switch t := field.Type.(type) {
			case *ast.Ident: // simple types e.g. int, string, etc.
				output[structName][fname] = FieldInfo{Type: t.Name}

			case *ast.SelectorExpr: // types like time.Time, time.Duration, nulls.String
				typ := fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
				output[structName][fname] = FieldInfo{Type: typ}

			case *ast.StarExpr: // pointer to something
				switch se := t.X.(type) {
				case *ast.Ident: // *SomeStruct, *string, *int etc.
					typ := se.Name
					output[structName][fname] = FieldInfo{Type: typ, IsPointer: true}
				case *ast.SelectorExpr: // *time.Time
					typ := fmt.Sprintf("%s.%s", se.X.(*ast.Ident).Name, se.Sel.Name)
					output[structName][fname] = FieldInfo{Type: typ, IsPointer: true}
				default:
					typ := fmt.Sprintf("%s", reflect.TypeOf(t))
					output[structName]["unsupported_star_expr_"+typ] = FieldInfo{Type: fmt.Sprintf("%T", se)}
					return true
				}

			case *ast.ArrayType:
				typ := "empty_type"
				switch at := t.Elt.(type) {
				case *ast.SelectorExpr:
					typ = at.Sel.Name
				case *ast.Ident:
					typ = at.Name
				default:
					typ := fmt.Sprintf("%s", reflect.TypeOf(t))
					output[structName]["unsupported_array_type_"+typ] = FieldInfo{Type: fmt.Sprintf("%T", at)}
					return true
				}
				output[structName][fname] = FieldInfo{Type: typ}

			default:
				typ := fmt.Sprintf("%s", reflect.TypeOf(t))
				output[structName]["unsupported_"+typ] = FieldInfo{Type: typ}
			}
		}
		return false
	}
}

// Parse gets path to source file or content of source file as a io.Reader and
// run inspect functions on it. Function returns list of structures with their
// fields.
func Parse(path string, src io.Reader) (StructureList, error) {
	node, err := parser.ParseFile(token.NewFileSet(), path, src, 0)
	if err != nil {
		return nil, err
	}

	info := StructureList{}

	ast.Inspect(node, inspect(info))

	return info, nil
}

// Lookup return structure by name from parsed source file or an error if
// structure with such name not found.
func Lookup(sl StructureList, structName string) (Structure, error) {
	f, ok := sl[structName]
	if !ok {
		return f, fmt.Errorf("structure %q not found", structName)
	}

	return f, nil
}
