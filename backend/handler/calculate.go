package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/mehmetkr/forge-abacus/backend/calc"
)

// Calculate handles POST /api/calculate requests.
func Calculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &maxBytesErr):
			writeError(w, http.StatusBadRequest, "request body too large")
		case errors.As(err, &typeErr) && (typeErr.Field == "a" || typeErr.Field == "b"):
			if strings.HasPrefix(typeErr.Value, "number") {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("operand '%s' is not a finite number", typeErr.Field))
			} else {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("field '%s' must be a number", typeErr.Field))
			}
		default:
			writeError(w, http.StatusBadRequest, "invalid JSON in request body")
		}
		return
	}

	if req.Operation == "" {
		writeError(w, http.StatusBadRequest, "field 'operation' is required")
		return
	}
	if req.A == nil {
		writeError(w, http.StatusBadRequest, "field 'a' is required")
		return
	}
	if req.Operation != "sqrt" && req.B == nil {
		writeError(w, http.StatusBadRequest, "field 'b' is required")
		return
	}
	if math.IsInf(*req.A, 0) || math.IsNaN(*req.A) {
		writeError(w, http.StatusBadRequest, "operand 'a' is not a finite number")
		return
	}
	if req.B != nil && (math.IsInf(*req.B, 0) || math.IsNaN(*req.B)) {
		writeError(w, http.StatusBadRequest, "operand 'b' is not a finite number")
		return
	}

	var b float64
	if req.B != nil {
		b = *req.B
	}
	result, err := calc.Compute(req.Operation, *req.A, b)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CalculateResponse{Result: result})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
