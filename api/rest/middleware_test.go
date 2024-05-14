//go:build unit
// +build unit

package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapMiddleware(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		middlewares []Middleware
		handler     Handler
		expected    string // Expected response body
	}{
		{
			name:        "No Middleware",
			middlewares: nil,
			handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				w.Write([]byte("original"))
				return nil
			},
			expected: "original",
		},
		{
			name: "Single Middleware",
			middlewares: []Middleware{
				func(next Handler) Handler {
					return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
						w.Write([]byte("before "))
						err := next(ctx, w, r)
						w.Write([]byte(" after"))
						return err
					}
				},
			},
			handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				w.Write([]byte("original"))
				return nil
			},
			expected: "before original after",
		},
		{
			name: "Multiple Middlewares",
			middlewares: []Middleware{
				func(next Handler) Handler {
					return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
						w.Write([]byte("outer1 "))
						err := next(ctx, w, r)
						w.Write([]byte(" outer2"))
						return err
					}
				},
				func(next Handler) Handler {
					return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
						w.Write([]byte("inner1 "))
						err := next(ctx, w, r)
						w.Write([]byte(" inner2"))
						return err
					}
				},
			},
			handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				w.Write([]byte("original"))
				return nil
			},
			expected: "outer1 inner1 original inner2 outer2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Wrap the handler with the given middlewares
			wrappedHandler := wrapMiddleware(tc.middlewares, tc.handler)

			// Create a test request
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			// Call the wrapped handler
			err := wrappedHandler(context.Background(), w, req)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check the response
			res := w.Result()
			if res.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: got %d, want %d", res.StatusCode, http.StatusOK)
			}
			body := w.Body.String()
			if body != tc.expected {
				t.Errorf("Unexpected body: got %q, want %q", body, tc.expected)
			}
		})
	}
}
