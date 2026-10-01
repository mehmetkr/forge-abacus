package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mehmetkr/forge-abacus/backend/handler"
)

func assertCORSHeaders(t *testing.T, h http.Header) {
	t.Helper()
	if v := h.Get("Access-Control-Allow-Origin"); v != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", v, "*")
	}
	if v := h.Get("Access-Control-Allow-Methods"); v != "POST, OPTIONS" {
		t.Fatalf("Access-Control-Allow-Methods = %q, want %q", v, "POST, OPTIONS")
	}
	if v := h.Get("Access-Control-Allow-Headers"); v != "Content-Type" {
		t.Fatalf("Access-Control-Allow-Headers = %q, want %q", v, "Content-Type")
	}
}

func TestRouting(t *testing.T) {
	h := newMux()

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantError  string
	}{
		{name: "GET returns 405", method: http.MethodGet, wantStatus: 405, wantError: "method not allowed"},
		{name: "PUT returns 405", method: http.MethodPut, wantStatus: 405, wantError: "method not allowed"},
		{name: "DELETE returns 405", method: http.MethodDelete, wantStatus: 405, wantError: "method not allowed"},
		{name: "OPTIONS returns 204", method: http.MethodOptions, wantStatus: 204},
		{name: "POST includes CORS headers", method: http.MethodPost, body: `{"operation":"add","a":1,"b":2}`, wantStatus: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}
			req := httptest.NewRequest(tt.method, "/api/calculate", body)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			assertCORSHeaders(t, rec.Header())

			if tt.wantError != "" {
				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
				}
				var resp handler.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if resp.Error != tt.wantError {
					t.Fatalf("error = %q, want %q", resp.Error, tt.wantError)
				}
			}
		})
	}
}

func TestEndToEnd(t *testing.T) {
	srv := httptest.NewServer(newMux())
	defer srv.Close()

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantResult float64
		wantError  string
	}{
		{name: "successful computation", method: http.MethodPost, body: `{"operation":"add","a":2,"b":3}`, wantStatus: 200, wantResult: 5},
		{name: "validation error", method: http.MethodPost, body: `{"operation":"add","a":1}`, wantStatus: 400, wantError: "field 'b' is required"},
		{name: "method not allowed", method: http.MethodGet, wantStatus: 405, wantError: "method not allowed"},
		{name: "CORS preflight", method: http.MethodOptions, wantStatus: 204},
		{name: "malformed body", method: http.MethodPost, body: `{invalid}`, wantStatus: 400, wantError: "invalid JSON in request body"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}
			req, err := http.NewRequest(tt.method, srv.URL+"/api/calculate", body)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			assertCORSHeaders(t, resp.Header)

			if tt.wantError != "" {
				var errResp handler.ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != tt.wantError {
					t.Fatalf("error = %q, want %q", errResp.Error, tt.wantError)
				}
			} else if tt.wantStatus == 200 {
				var successResp handler.CalculateResponse
				if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
					t.Fatalf("failed to decode success response: %v", err)
				}
				if successResp.Result != tt.wantResult {
					t.Fatalf("result = %v, want %v", successResp.Result, tt.wantResult)
				}
			}
		})
	}
}
