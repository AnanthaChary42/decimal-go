package decimal

import "fmt"

// DecimalError represents an error during decimal operations.
type DecimalError struct {
	Message string
}

func (e *DecimalError) Error() string {
	return fmt.Sprintf("[DecimalError] %s", e.Message)
}

// Sentinel errors matching decimal.js error conditions.
var (
	ErrInvalidArg     = &DecimalError{Message: "Invalid argument"}
	ErrDivisionByZero = &DecimalError{Message: "Division by zero"}
	ErrPrecisionLimit = &DecimalError{Message: "Precision limit exceeded"}
)

func newInvalidArgError(val interface{}) error {
	return &DecimalError{Message: fmt.Sprintf("Invalid argument: %v", val)}
}
