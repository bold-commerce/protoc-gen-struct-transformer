package generator

import (
	"errors"
	"fmt"
)

var (
	// ErrNilOptions appears when object has no options.
	//
	// Options defined as
	//
	// messages MessageName {
	//	option (package.message_option_name) = value;
	//	int64 first_field = 1 [(package.field_option_name = value)];
	// }
	ErrNilOptions = errors.New("options are nil")

	// ErrFileSkipped is returned when .proto file has not go_models_file_path
	// option.
	ErrFileSkipped = errors.New("files was skipped")
)

// errOptionNotExists represent option extract-related error.
type errOptionNotExists string

// newErrOptionNotExists initializes new option-related error.
func newErrOptionNotExists(optName string) error {
	return errOptionNotExists(fmt.Sprintf("option %q does not exists", optName))
}

// Error is an error interface implementation.
func (e errOptionNotExists) Error() string {
	return string(e)
}

// loggableError is an error type which should be added to result file as a
// comment if happened.
type loggableError struct {
	message string
}

// newLoggableError initializes error which will be added to output file.
func newLoggableError(format string, args ...interface{}) loggableError {
	return loggableError{message: fmt.Sprintf(format, args...)}
}

// Error is an error interface implementation.
func (e loggableError) Error() string {
	return e.message
}
