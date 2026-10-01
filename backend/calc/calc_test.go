package calc

import (
	"errors"
	"math"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		a, b     float64
		want     float64
		wantErr  error
	}{
		// Addition
		{name: "add two positives", operator: "add", a: 2, b: 3, want: 5},
		{name: "add negative numbers", operator: "add", a: -2, b: -3, want: -5},
		{name: "add zero", operator: "add", a: 5, b: 0, want: 5},

		// Subtraction
		{name: "subtract two positives", operator: "subtract", a: 5, b: 3, want: 2},
		{name: "subtract resulting negative", operator: "subtract", a: 3, b: 5, want: -2},

		// Multiplication
		{name: "multiply two positives", operator: "multiply", a: 4, b: 3, want: 12},
		{name: "multiply by zero", operator: "multiply", a: 5, b: 0, want: 0},
		{name: "multiply negatives", operator: "multiply", a: -3, b: -4, want: 12},

		// Division
		{name: "divide evenly", operator: "divide", a: 10, b: 2, want: 5},
		{name: "divide with remainder", operator: "divide", a: 10, b: 3, want: 10.0 / 3.0},
		{name: "divide negative", operator: "divide", a: -10, b: 2, want: -5},

		// Power
		{name: "power of two", operator: "power", a: 2, b: 10, want: 1024},
		{name: "power of one", operator: "power", a: 5, b: 1, want: 5},
		{name: "negative exponent", operator: "power", a: 2, b: -1, want: 0.5},

		// Square root
		{name: "sqrt of perfect square", operator: "sqrt", a: 9, b: 0, want: 3},
		{name: "sqrt of two", operator: "sqrt", a: 2, b: 0, want: math.Sqrt(2)},
		{name: "sqrt of zero", operator: "sqrt", a: 0, b: 0, want: 0},

		// Percentage
		{name: "percentage basic", operator: "percentage", a: 15, b: 200, want: (15.0 / 100.0) * 200.0},
		{name: "percentage full", operator: "percentage", a: 100, b: 50, want: 50},
		{name: "percentage zero", operator: "percentage", a: 0, b: 200, want: 0},

		// Edge cases: division by zero
		{name: "divide by zero", operator: "divide", a: 1, b: 0, wantErr: ErrDivisionByZero},
		{name: "zero divided by zero", operator: "divide", a: 0, b: 0, wantErr: ErrDivisionByZero},

		// Edge cases: negative square root
		{name: "sqrt of negative", operator: "sqrt", a: -1, b: 0, wantErr: ErrNegativeSqrt},

		// Edge cases: 0^0 returns 1 (IEEE 754)
		{name: "zero to the power of zero", operator: "power", a: 0, b: 0, want: 1},

		// Edge cases: overflow to infinity
		{name: "overflow to infinity", operator: "power", a: 1e308, b: 2, wantErr: ErrResultInfinite},
		{name: "addition overflow", operator: "add", a: math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrResultInfinite},

		// Edge cases: post-computation catch-all with non-finite inputs
		{name: "inf input produces inf result", operator: "add", a: math.Inf(1), b: 1, wantErr: ErrResultInfinite},
		{name: "nan input produces nan result", operator: "add", a: math.NaN(), b: 1, wantErr: ErrResultUndefined},

		// Edge cases: finite inputs producing NaN
		{name: "power produces NaN", operator: "power", a: -1, b: 0.5, wantErr: ErrResultUndefined},

		// Edge cases: unknown operator
		{name: "unsupported operation", operator: "modulo", a: 5, b: 3, wantErr: ErrUnsupportedOperation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compute(tt.operator, tt.a, tt.b)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
