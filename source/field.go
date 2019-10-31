package source

import "fmt"

type (
	// FieldInfo containt information about one structure field without field name.
	FieldInfo struct {
		// Field type name.
		Type string
		// Equals true if field is a pointer.
		IsPointer bool
	}

	// Structure is a set of fields of one structure.
	Structure map[string]FieldInfo
	// StructureList is a list of parsed structures.
	StructureList map[string]Structure
)

// String return structure information as a string.
func (s Structure) String() string {
	c := "\n// Target struct fields:\n"
	for k, v := range s {
		c += fmt.Sprintf("// Field: %q, Type: %q, isPointer: %t\n", k, v.Type, v.IsPointer)
	}
	c += "\n"
	return c
}
