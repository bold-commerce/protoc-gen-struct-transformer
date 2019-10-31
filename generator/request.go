package generator

import (
	"strings"
)

// Setter is a interface which allows to set map key with value. Both are of
// type string.
type Setter interface {
	Set(string, string) error
}

// SetParameters accepts flag.CommandLine as a setter, string with params
// from protobuf compile input. Functions decodes params and adds it as command
// line flags.
func SetParameters(setter Setter, param *string) error {
	if param == nil {
		return nil
	}

	for _, p := range strings.Split(*param, ",") {
		spec := strings.SplitN(p, "=", 2)
		// skip output dir and modifiers
		if len(spec) == 1 || strings.HasPrefix(spec[0], "M") {
			continue
		}

		if err := setter.Set(spec[0], spec[1]); err != nil {
			return err
		}
	}

	return nil
}
