package generator

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func xEntry(g2p, p2g string, gg, pp, swapped bool, expected string) TableEntry {
	args := []interface{}{g2p, p2g, gg, pp, swapped, expected}

	desc := fmt.Sprintf("GoToProtoType: %q, ProtoToGoType: %q, GoIsPointer: %t, ProtoIsPointer: %t, swapped: %t, expected: %q", args...)
	return Entry(desc, args...)
}

var _ = Describe("Template", func() {

	Describe("Field methods", func() {

		var field = &Field{
			Name:      "name",
			ProtoName: "proto_name",
		}

		Context("when call name(swapped) method", func() {

			It("returns field.Name if swapped==true", func() {
				name := field.name(true)
				Expect(name).To(Equal("name"))
			})

			It("returns field.Name if swapped==false", func() {
				name := field.name(false)
				Expect(name).To(Equal("proto_name"))
			})
		})

		Context("when call convertFunc(swapped) method", func() {

			DescribeTable("check result",
				// func(swapped bool, expected string, f *Field) {
				func(g2p, p2g string, gg, pp, swapped bool, expected string) {
					f := &Field{
						GoToProtoType:  g2p,
						ProtoToGoType:  p2g,
						GoIsPointer:    gg,
						ProtoIsPointer: pp,
					}

					r := f.convertFunc(swapped)
					Expect(r).To(Equal(expected))
				},

				xEntry("go2proto", "proto2go", false, false, false, "proto2go"),
				xEntry("go2proto", "proto2go", false, false, true, "go2proto"),
				xEntry("go2proto", "proto2go", true, true, false, "proto2goPtr"),
				xEntry("go2proto", "proto2go", true, true, true, "go2protoPtr"),
				xEntry("go2proto", "proto2go", false, true, false, "proto2goPtrVal"),
				xEntry("go2proto", "proto2go", false, true, true, "go2protoValPtr"),
				xEntry("go2proto", "proto2go", true, false, false, "proto2goValPtr"),
				xEntry("go2proto", "proto2go", true, false, true, "go2protoPtrVal"),
				xEntry("go2protoList", "proto2goList", false, false, false, "proto2go"),
				xEntry("go2protoList", "proto2goList", false, false, true, "go2proto"),
				xEntry("go2protoList", "proto2goList", true, true, false, "proto2goListPtr"),
				xEntry("go2protoList", "proto2goList", true, true, true, "go2protoListPtr"),
				xEntry("go2protoList", "proto2goList", true, false, false, "proto2goValPtrList"),
				xEntry("go2protoList", "proto2goList", true, false, true, "go2protoPtrValList"),
				xEntry("go2protoList", "proto2goList", false, true, false, "proto2goPtrValList"),
				xEntry("go2protoList", "proto2goList", false, true, true, "go2protoValPtrList"),
			)
		})
	})

	Describe("formatOneofField", func() {

		DescribeTable("check returns",
			func(f Field, swapped bool, pref, expected string) {
				r := formatOneofField(f, swapped, pref)
				Expect(r).To(Equal(expected))
			},

			Entry("Field has no oneof declaration", Field{
				Name:      "field_name",
				OneofDecl: "",
			}, false, "", `/* field "field_name" is not Oneof field*/`),

			Entry("ProtoToGoType is empty", Field{
				ProtoName: "proto_name",
				OneofDecl: "oneof_decl_name",
			}, false, "", "src.proto_name"),

			Entry("ProtoToGoType is empty, swapped", Field{
				ProtoType: "proto_type",
				OneofDecl: "oneof_decl_name",
			}, true, "prefix", "&prefix.proto_type{}"),

			Entry("ProtoToGoType is not empty", Field{
				ProtoName:     "proto_name",
				ProtoToGoType: "p2g",
				OneofDecl:     "oneof_decl_name",
			}, false, "prefix", "p2g(src.proto_name)"),
		)
	})

	Describe("formatOneofInitField", func() {

		DescribeTable("check returns",
			func(f Field, swapped bool, expected string) {
				r := formatOneofInitField(f, swapped)
				Expect(r).To(Equal(expected))
			},
			Entry("Not oneof", Field{}, false, ""),
			Entry("Oneof field", Field{OneofDecl: "decl"}, false, ""),
			Entry("Not oneof, swapped", Field{}, true, ""),
			Entry("", Field{
				OneofDecl:     "oneof_decl_name",
				GoToProtoType: "g2p",
				Name:          "field_name",
				ProtoName:     "proto_name",
			}, true, " g2p(src.field_name, s.proto_name, version)"),
		)
	})

	Describe("formatComplexField", func() {

		DescribeTable("check returns",
			func(f Field, swapped bool, expected string) {
				r := formatComplexField(f, swapped)
				Expect(r).To(Equal(expected))
			},
			Entry("ProtoToGoType is not empty", Field{
				Name:           "name",
				ProtoName:      "proto_name",
				GoToProtoType:  "g2p",
				ProtoToGoType:  "p2g",
				GoIsPointer:    true,
				ProtoIsPointer: false,
				Opts:           ", opts...",
			}, false, " p2gValPtr(src.proto_name , opts...)"),

			Entry("ProtoToGoType is empty", Field{
				Name:           "name",
				ProtoName:      "proto_name",
				GoToProtoType:  "g2p",
				ProtoToGoType:  "",
				GoIsPointer:    true,
				ProtoIsPointer: false,
				Opts:           ", opts...",
			}, false, "src.proto_name"),
		)
	})

	Describe("formatField", func() {

		DescribeTable("check returns",
			func(f Field, swapped bool, pref, expected string) {
				r := formatField(f, swapped, pref)
				Expect(r).To(Equal(expected))
			},

			Entry("Not oneof", Field{
				Name:      "name",
				ProtoName: "proto_name",
			}, false, "", "name: src.proto_name,"),

			Entry("Not oneof, swapped", Field{
				Name:      "name",
				ProtoName: "proto_name",
			}, true, "", "proto_name: src.name,"),

			Entry("Oneof", Field{
				Name:      "name",
				ProtoName: "proto_name",
				OneofDecl: "oneof_decl_name",
			}, false, "prefix", "name: src.proto_name,"),

			Entry("Oneof, swapped", Field{
				Name:      "name",
				ProtoName: "proto_name",
				ProtoType: "proto_type",
				OneofDecl: "oneof_decl_name",
			}, true, "prefix", "proto_name: &prefix.proto_type{},"),
		)
	})

	Describe("Data.Swap", func() {

		Context("when Swap() called", func() {

			var d *Data

			BeforeEach(func() {
				d = &Data{
					Src:        "src",
					Dst:        "dst",
					SrcPref:    "src_pref",
					DstPref:    "dst_pref",
					SrcFn:      "src_fn",
					DstFn:      "dst_fn",
					SrcPointer: "src_pointer",
					DstPointer: "dst_pointer",
					Swapped:    false,
				}
			})

			It("swap some fields", func() {
				d.swap()
				Expect(*d).To(MatchFields(IgnoreExtras, Fields{
					"Src":        Equal("dst"),
					"Dst":        Equal("src"),
					"SrcPref":    Equal("dst_pref"),
					"DstPref":    Equal("src_pref"),
					"SrcFn":      Equal("dst_fn"),
					"DstFn":      Equal("src_fn"),
					"SrcPointer": Equal("dst_pointer"),
					"DstPointer": Equal("src_pointer"),
					"Swapped":    BeTrue(),
				}))
			})
		})
	})

	Describe("templateWithHelpers", func() {
		var w *bytes.Buffer

		BeforeEach(func() {
			w = bytes.NewBuffer([]byte{})
		})

		Context("when execute whole template", func() {

			It("returns full function set as string", func() {
				t, err := templateWithHelpers("test_template")
				Expect(err).NotTo(HaveOccurred())

				err = t.Execute(w, Data{
					SrcPref:       "src_pref",
					Src:           "src",
					SrcFn:         "src_fn",
					SrcPointer:    "src_pointer",
					DstPref:       "dst_pref",
					Dst:           "dst",
					DstFn:         "dst_fn",
					DstPointer:    "dst_pointer",
					Swapped:       false,
					HelperPackage: "hp",
					Ptr:           false,
					Fields: []Field{
						{
							Name:           "FirstField",
							ProtoName:      "proto_name",
							ProtoType:      "proto_type",
							ProtoToGoType:  "FirstProto2go",
							GoToProtoType:  "FirstGo2proto",
							GoIsPointer:    false,
							ProtoIsPointer: false,
							UsePackage:     false,
							OneofDecl:      "",
							Opts:           "",
						},
						{
							Name:           "SecondField",
							ProtoName:      "proto_name2",
							ProtoType:      "proto_type2",
							ProtoToGoType:  "SecondProto2go",
							GoToProtoType:  "SecondGo2proto",
							GoIsPointer:    false,
							ProtoIsPointer: false,
							UsePackage:     false,
							OneofDecl:      "oneof_decl_name",
							Opts:           "",
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			})

		})

	})

	Describe("Template parts", func() {
		var w *bytes.Buffer

		BeforeEach(func() {
			w = bytes.NewBuffer([]byte{})
		})

		Context("when execute template funcNameT", func() {

			It("return formatted string", func() {
				d := Data{SrcFn: "SrcFn", DstFn: "DstFn"}
				err := funcNameT.Execute(w, d)
				Expect(err).NotTo(HaveOccurred())
				Expect(w.String()).To(Equal("SrcFnToDstFn"))
			})

		})

		Context("when execute template srcParamT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := srcParamT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Without prefix", Data{Src: "Src"}, "Src, opts ...TransformParam"),
				Entry("Withprefix", Data{SrcPref: "pref", Src: "Src"}, "pref.Src, opts ...TransformParam"),
			)

		})

		Context("when execute template dstParamT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := dstParamT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Without prefix", Data{Dst: "Dst"}, "Dst"),
				Entry("With prefix", Data{DstPref: "pref", Dst: "Dst"}, "pref.Dst"),
			)

		})

		Context("when execute template ptrValT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptrValT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Not swapped", Data{Swapped: false}, "PtrVal"),
				Entry("Swapped", Data{Swapped: true}, "ValPtr"),
			)
		})

		Context("when execute template ptrT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptrT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{Ptr: true}, "Ptr"),
				Entry("Not Ptr", Data{Ptr: false}, "Val"),
			)
		})

		Context("when execute template ptrOnlyT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptrOnlyT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{Ptr: true}, "Ptr"),
				Entry("Not Ptr", Data{Ptr: false}, ""),
			)
		})

		Context("when execute template starT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := starT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{Ptr: true}, "*"),
				Entry("Not Ptr", Data{Ptr: false}, ""),
			)
		})

		Context("when execute template ptr2ptrT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptr2ptrT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnPtr(src *SrcPref.Src, opts ...TransformParam) *DstPref.Dst {
	if src == nil {
		return nil
	}

	d := SrcFnToDstFn(*src, opts...)
	return &d
}`),
			)
		})

		Context("when execute template ptr2valT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptr2valT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnPtrVal(src *SrcPref.Src, opts ...TransformParam) DstPref.Dst {
	if src == nil {
		return DstPref.Dst{}
	}

	return SrcFnToDstFn(*src, opts...)
}`),
			)
		})

		Context("when execute template val2ptrT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := val2ptrT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnValPtr(src SrcPref.Src, opts ...TransformParam) *DstPref.Dst {
	d := SrcFnToDstFn(src, opts...)
	return &d
}`),
			)
		})

		Context("when execute template val2valT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := val2valT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
					Fields: []Field{
						{
							Name:           "FirstField",
							ProtoName:      "proto_name",
							ProtoType:      "proto_type",
							ProtoToGoType:  "FirstGo2proto",
							GoToProtoType:  "FirstProto2go",
							GoIsPointer:    false,
							ProtoIsPointer: false,
							UsePackage:     false,
							OneofDecl:      "",
							Opts:           "",
						},
						{
							Name:           "SecondField",
							ProtoName:      "proto_name2",
							ProtoType:      "proto_type2",
							ProtoToGoType:  "SecondProto2go",
							GoToProtoType:  "SecondGo2proto",
							GoIsPointer:    false,
							ProtoIsPointer: false,
							UsePackage:     false,
							OneofDecl:      "oneof_decl_name",
							Opts:           "",
						},
					},
				}, `func SrcFnToDstFn(src SrcPref.Src, opts ...TransformParam) DstPref.Dst {
	s := DstPref.Dst{
			FirstField:  FirstGo2proto(src.proto_name ),
			SecondField: SecondProto2go(src.proto_name2),
	}

	applyOptions(opts...)



	return s
}`),
			)
		})

		Context("when execute template lst2lstT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := lst2lstT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnValList(src []SrcPref.Src, opts ...TransformParam) []DstPref.Dst {
	resp := make([]DstPref.Dst, len(src))

	for i, s := range src {
		resp[i] = SrcFnToDstFn(s, opts...)
	}

	return resp
}`),
			)
		})

		Context("when execute template ptrlst2ptrlstT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptrlst2ptrlstT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnPtrList(src []*SrcPref.Src, opts ...TransformParam) []*DstPref.Dst {
	resp := make([]*DstPref.Dst, len(src))

	for i, s := range src {
		resp[i] = SrcFnToDstFnPtr(s, opts...)
	}

	return resp
}`),
			)
		})

		Context("when execute template vallst2vallstT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := vallst2vallstT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnValList(src []SrcPref.Src, opts ...TransformParam) []DstPref.Dst {
	resp := make([]DstPref.Dst, len(src))

	for i, s := range src {
		resp[i] = SrcFnToDstFn(s, opts...)
	}

	return resp
}`),
			)
		})

		Context("when execute template ptrlst2vallstT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptrlst2vallstT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `func SrcFnToDstFnPtrValList(src []SrcPref.Src, opts ...TransformParam) []DstPref.Dst {
	resp := make([]DstPref.Dst, len(src))

	for i, s := range src {
		resp[i] = SrcFnToDstFn(*s)
		}

	return resp
}`),
			)
		})

		Context("when execute template ptr2vallstT", func() {

			DescribeTable("check result",
				func(d Data, expected string) {
					err := ptr2vallstT.Execute(w, d)
					Expect(err).NotTo(HaveOccurred())
					Expect(w.String()).To(Equal(expected))
				},
				Entry("Ptr", Data{
					Src:     "Src",
					SrcFn:   "SrcFn",
					SrcPref: "SrcPref",
					Dst:     "Dst",
					DstFn:   "DstFn",
					DstPref: "DstPref",
				}, `// SrcFnToDstFnList is DEPRECATED. Use SrcFnToDstFnPtrValList instead.
func SrcFnToDstFnList(src []SrcPref.Src, opts ...TransformParam) []DstPref.Dst {
	return SrcFnToDstFnPtrValList(src)
}`),
			)
		})
	})

})
