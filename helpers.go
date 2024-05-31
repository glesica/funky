package funky

import "errors"

var (
	// Stop indicates that iteration should not continue. The error
	// will generally be bubbled up to the caller.
	Stop = errors.New("stop iteration")
)

// Done is a helper that can be used to indicate that an
// iterator has finished producing values, whether because
// of an error, or otherwise.
//
// todo: this needn't take an error now that we allow errors mid-iteration
func Done[T any](err error) (T, bool, error) {
	return *new(T), false, err
}

func Value[T any](val T, err error) (T, bool, error) {
	return val, true, err
}
