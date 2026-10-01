package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// All tests use httptest.NewRequest which does not set Content-Type,
// implicitly verifying the handler succeeds without a Content-Type header.

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult float64
		wantError  string
	}{
		// Success: one per operation
		{name: "add", body: `{"operation":"add","a":2,"b":3}`, wantStatus: 200, wantResult: 5},
		{name: "subtract", body: `{"operation":"subtract","a":5,"b":3}`, wantStatus: 200, wantResult: 2},
		{name: "multiply", body: `{"operation":"multiply","a":4,"b":3}`, wantStatus: 200, wantResult: 12},
		{name: "divide", body: `{"operation":"divide","a":10,"b":3}`, wantStatus: 200, wantResult: 10.0 / 3.0},
		{name: "power", body: `{"operation":"power","a":2,"b":10}`, wantStatus: 200, wantResult: 1024},
		{name: "sqrt", body: `{"operation":"sqrt","a":9}`, wantStatus: 200, wantResult: 3},
		{name: "percentage", body: `{"operation":"percentage","a":15,"b":200}`, wantStatus: 200, wantResult: (15.0 / 100.0) * 200.0},

		// Validation: missing fields
		{name: "missing operation", body: `{"a":1,"b":2}`, wantStatus: 400, wantError: "field 'operation' is required"},
		{name: "missing a", body: `{"operation":"add","b":2}`, wantStatus: 400, wantError: "field 'a' is required"},
		{name: "missing b", body: `{"operation":"add","a":1}`, wantStatus: 400, wantError: "field 'b' is required"},

		// Validation: wrong types
		{name: "wrong type for a", body: `{"operation":"add","a":"hello","b":2}`, wantStatus: 400, wantError: "field 'a' must be a number"},
		{name: "wrong type for b", body: `{"operation":"add","a":1,"b":true}`, wantStatus: 400, wantError: "field 'b' must be a number"},

		// Validation: non-finite input
		{name: "infinite a", body: `{"operation":"add","a":1e309,"b":1}`, wantStatus: 400, wantError: "operand 'a' is not a finite number"},
		{name: "infinite b", body: `{"operation":"add","a":1,"b":1e309}`, wantStatus: 400, wantError: "operand 'b' is not a finite number"},

		// Computation errors
		{name: "divide by zero", body: `{"operation":"divide","a":1,"b":0}`, wantStatus: 400, wantError: "division by zero is undefined"},
		{name: "negative sqrt", body: `{"operation":"sqrt","a":-1}`, wantStatus: 400, wantError: "square root of negative number is undefined"},

		// Body errors
		{name: "malformed JSON", body: `{invalid}`, wantStatus: 400, wantError: "invalid JSON in request body"},
		{name: "empty body", body: ``, wantStatus: 400, wantError: "invalid JSON in request body"},
		{name: "oversized body", body: `{"a":` + strings.Repeat("1", 1<<20) + `}`, wantStatus: 400, wantError: "request body too large"},

		// Special cases
		{name: "unknown operation", body: `{"operation":"modulo","a":5,"b":3}`, wantStatus: 400, wantError: "unsupported operation: modulo"},
		{name: "unknown fields ignored", body: `{"operation":"add","a":2,"b":3,"extra":"field"}`, wantStatus: 200, wantResult: 5},
		{name: "sqrt ignores b", body: `{"operation":"sqrt","a":9,"b":99}`, wantStatus: 200, wantResult: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			Calculate(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
			}

			if tt.wantError != "" {
				var resp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if resp.Error != tt.wantError {
					t.Fatalf("error = %q, want %q", resp.Error, tt.wantError)
				}
			} else {
				var resp CalculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode success response: %v", err)
				}
				if resp.Result != tt.wantResult {
					t.Fatalf("result = %v, want %v", resp.Result, tt.wantResult)
				}
			}
		})
	}
}
