package lwcomponent

import "github.com/pkg/errors"

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
