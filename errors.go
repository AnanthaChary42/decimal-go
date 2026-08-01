package decimal

import "fmt"

// DecimalError is the error type returned by decimal operations.
type DecimalError struct {
	Message string
}

func (e *DecimalError) Error() string {
	return "[DecimalError] " + e.Message
}

// newError creates a new DecimalError.
func newError(msg string) *DecimalError {
	return &DecimalError{Message: msg}
}

// newInvalidArgError creates an "Invalid argument" error.
func newInvalidArgError(v interface{}) *DecimalError {
	return &DecimalError{Message: fmt.Sprintf("Invalid argument: %v", v)}
}

// Sentinel errors.
var (
	ErrPrecisionLimit = &DecimalError{Message: "Precision limit exceeded"}
)
