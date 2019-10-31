package generator

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// setter is a faked setter.
type setter map[string]string

func (s setter) Set(k, v string) error {
	if k == "err" {
		return errors.New("key_error")
	}
	s[k] = v

	return nil
}

func (s setter) String() string {
	r := ""
	for k, v := range s {
		r += fmt.Sprintf("%s=%s,", k, v)
	}

	return r
}

var _ = Describe("Request", func() {

	Describe("SetParameters", func() {

		var set = setter{}

		DescribeTable("check results",
			func(s Setter, p *string, expected string) {
				err := SetParameters(s, p)
				Expect(err).NotTo(HaveOccurred())

				Expect(s.(fmt.Stringer).String()).To(Equal(expected))
			},
			Entry("nil as param", set, nil, ""),
			Entry("", set, sp("Mkey1=val1"), ""),
			Entry("", set, sp("Mkey1=val1,Mkey2=val2"), ""),
			Entry("", set, sp("Mkey1=val1,Mkey2=val2,key3=val3"), "key3=val3,"),
		)

	})
})
