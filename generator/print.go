package generator

import (
	"fmt"
	"io"
)

// p is a debug function which writes formatted sting into w if w is not nil.
func p(w io.Writer, format string, a ...interface{}) {
	if w == nil {
		return
	}

	fmt.Fprintf(w, format, a...)
}
