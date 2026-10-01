package calc

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrDivisionByZero       = errors.New("division by zero is undefined")
	ErrNegativeSqrt         = errors.New("square root of negative number is undefined")
	ErrResultInfinite       = errors.New("result is infinite")
	ErrResultUndefined      = errors.New("result is undefined")
	ErrUnsupportedOperation = errors.New("unsupported operation")
)

// Compute performs the given arithmetic operation. Operands must be finite.
func Compute(operator string, a, b float64) (float64, error) {
	var result float64

	switch operator {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result = a / b
	case "power":
		result = math.Pow(a, b)
	case "sqrt":
		if a < 0 {
			return 0, ErrNegativeSqrt
		}
		result = math.Sqrt(a)
	case "percentage":
		result = (a / 100) * b
	default:
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedOperation, operator)
	}

	if math.IsInf(result, 0) {
		return 0, ErrResultInfinite
	}
	if math.IsNaN(result) {
		return 0, ErrResultUndefined
	}

	return result, nil
}
