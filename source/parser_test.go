package source

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {

	DescribeTable("check result",
		func(fileContent string, expected StructureList) {
			str, err := Parse("file.go", bytes.NewReader([]byte(fileContent)))
			Expect(err).NotTo(HaveOccurred())

			Expect(str).To(Equal(expected))
		},

		Entry("File without structures", `package model`, StructureList{}),

		Entry("File with non-struct types", `package model

type myInt int
type myString string
`, StructureList{}),

		Entry("File with one struct", `package model

type (
	MyStruct struct {
		ID       int
		Name     string
	}
)`, StructureList{
			"MyStruct": {
				"ID":   {Type: "int", IsPointer: false},
				"Name": {Type: "string", IsPointer: false},
			},
		}),

		Entry("File with two structs, one is embedded into another", `package model

type (
	Comment struct {
		Content string
	}

	MyStruct struct {
		ID       int
		Name     string
		Comment
	}
)`, StructureList{
			"Comment": {
				"Content": {Type: "string", IsPointer: false},
			},
			"MyStruct": {
				"embedded_0": {Type: "Comment", IsPointer: false},
				"Name":       {Type: "string", IsPointer: false},
				"ID":         {Type: "int", IsPointer: false},
			},
		}),

		Entry("File with two independent structs", `package model

type (
	Comment struct {
		ID			int
		Content string
	}

	MyStruct struct {
		ID       int
		Name     string
	}
)`, StructureList{
			"Comment": {
				"ID":      {Type: "int", IsPointer: false},
				"Content": {Type: "string", IsPointer: false},
			},
			"MyStruct": {
				"Name": {Type: "string", IsPointer: false},
				"ID":   {Type: "int", IsPointer: false},
			},
		}),

		Entry("File with one struct, fields are of type SelectorExpr: time.Time, etc.", `package model

type (
	MyStruct struct {
		ID				int
		Name			string
		CreatedAt time.Time
	}
)`, StructureList{
			"MyStruct": {
				"ID":        {Type: "int", IsPointer: false},
				"Name":      {Type: "string", IsPointer: false},
				"CreatedAt": {Type: "time.Time", IsPointer: false},
			},
		}),

		Entry("File with one struct, fields are of slice type.", `package model

type (
	MyStruct struct {
		ID					int
		Name				string
		SubMyStructs	[]int
	}
)`, StructureList{
			"MyStruct": {
				"ID":           {Type: "int", IsPointer: false},
				"Name":         {Type: "string", IsPointer: false},
				"SubMyStructs": {Type: "int", IsPointer: false},
			},
		}),

		Entry("File with one struct, fields are of struct slice type.", `package model

type (
	MyStruct struct {
		ID	 int
		Name string
		Tags []nulls.String
	}
)`, StructureList{
			"MyStruct": {
				"ID":   {Type: "int", IsPointer: false},
				"Name": {Type: "string", IsPointer: false},
				"Tags": {Type: "String", IsPointer: false},
			},
		}),

		Entry("File with one struct, fields are of unsupported slice type.", `package model

type (
	MyStruct struct {
		ID	 int
		Name string
		Items []map[string]int
	}
)`, StructureList{
			"MyStruct": {
				"ID":                     {Type: "int", IsPointer: false},
				"Name":                   {Type: "string", IsPointer: false},
				"unsupported_array_type": {Type: "empty_type", IsPointer: false},
			},
		}),

		Entry("File with one struct, field is of pointer type.", `package model

type (
	MyStruct struct {
		ID	 int
		Name *string
	}
)`, StructureList{
			"MyStruct": {
				"ID":   {Type: "int", IsPointer: false},
				"Name": {Type: "string", IsPointer: true},
			},
		}),

		Entry("File with one struct, field is of unsupported type.", `package model

type (
	MyStruct struct {
		ID	 int
		Name string
		F func()
		M map[int]string
	}
)`, StructureList{
			"MyStruct": {
				"ID":                        {Type: "int", IsPointer: false},
				"Name":                      {Type: "string", IsPointer: false},
				"unsupported_*ast.FuncType": {Type: "*ast.FuncType", IsPointer: false},
				"unsupported_*ast.MapType":  {Type: "*ast.MapType", IsPointer: false},
			},
		}),
	)

	Describe("Lookup", func() {

		Context("when call Lookup with existing struct", func() {

			It("returns set of fields", func() {
				str, err := Parse("file.go", bytes.NewReader([]byte(`package model

type (
	MyStruct struct {
		ID	 int
		Name *string
	}
)`)))
				Expect(err).NotTo(HaveOccurred())

				fields, err := Lookup(str, "MyStruct")
				Expect(err).NotTo(HaveOccurred())
				Expect(fields).To(Equal(Structure{
					"ID":   {Type: "int", IsPointer: false},
					"Name": {Type: "string", IsPointer: true},
				}))

			})
		})

		Context("when call Lookup with non-existing struct", func() {

			It("returns set of fields", func() {
				str, err := Parse("file.go", bytes.NewReader([]byte(`package model

type (
	MyStruct struct {
		ID	 int
		Name *string
	}
)`)))
				Expect(err).NotTo(HaveOccurred())

				fields, err := Lookup(str, "NotExists")
				Expect(err).To(MatchError(`structure "NotExists" not found`))
				Expect(fields).To(BeNil())
			})
		})
	})

})
