package lwcomponent

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrComponentNotFound = errors.New("component not found on disk")
)

// IsNotFound returns a boolean indicating whether the error is known to
// have determined the component is not found. It is satisfied by
// ErrNotApplyComment
//
func IsNotFound(err error) bool {
	return errors.Is(err, ErrComponentNotFound)
}

// RunError is a struct used to pass an error when a component tries to run and
// it fails, a few functions will return this error so that callers (upstream
// packages) can unwrap and identify that the error comes from this package
type RunError struct {
	ExitCode int
	Message  string
	Err      error
}

func (e *RunError) Error() string {
	if e.ExitCode == 0 {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
}

func (e *RunError) Unwrap() error {
	return e.Err
}
