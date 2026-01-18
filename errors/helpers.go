package errors

import "errors"

// As wraps errors.As, which is used for type assertion
func As(err error, target any) bool {
	return errors.As(err, &target)
}

// Is wraps errors.Is, which is used for error comparison
func Is(err, target error) bool {
	return errors.Is(err, target)
}
