//go:build unit
// +build unit

package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
)

// mockMiddleware is a Middleware implementation for testing purposes.
func mockMiddleware(next Handler) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Add test-specific logic here if needed.
		return next(ctx, w, r)
	}
}

// mockHandler is a Handler implementation for testing purposes.
func mockHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data":"some data", "success":true}`))
	return nil
}

// mockHandler is a Handler implementation for testing purposes.
func mockErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"errors":"some error"}`))
	return nil
}

func TestNew(t *testing.T) {
	shutdown := make(chan os.Signal, 1)

	api := New(shutdown, mockMiddleware)
	if api == nil {
		t.Fatal("New returned nil")
	}
	if api.shutdown != shutdown {
		t.Error("shutdown channel not assigned correctly")
	}
	if api.mux == nil {
		t.Error("mux not initialized")
	}
	if len(api.mw) != 1 {
		t.Error("middleware not added correctly")
	}
}

func TestSignalShutdown(t *testing.T) {
	shutdown := make(chan os.Signal, 1)
	api := New(shutdown)

	// Start a goroutine to capture the signal
	go func() {
		time.Sleep(2 * time.Second)
		api.SignalShutdown()
	}()

	select {
	case sig := <-shutdown:
		if sig != syscall.SIGTERM {
			t.Errorf("Expected signal %v, got %v", syscall.SIGTERM, sig)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("TestSignalShutdown timed out")
	}
}

func TestAPI_Handle(t *testing.T) {
	shutdown := make(chan os.Signal, 1)
	api := New(shutdown)

	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		api.Handle("GET", "/test", mockHandler)
		api.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}
		body := rr.Body.String()
		response := Response{}
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Fatal("could not unmarshal response body")
		}
		if response.Success != true || response.Data == nil || response.Errors != nil {
			t.Errorf("handler returned wrong response: got %v", body)
		}
	})

	t.Run("Error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/error", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		api.Handle("GET", "/error", mockErrorHandler)
		api.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}
		body := rr.Body.String()
		response := Response{}
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Fatal("could not unmarshal response body")
		}
		if response.Success != false || response.Data != nil || response.Errors != "some error" {
			t.Errorf("handler returned wrong response: got %v", body)
		}
	})

	t.Run("Middleware", func(t *testing.T) {
		middlewareCalled := false
		testMiddleware := func(next Handler) Handler {
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				middlewareCalled = true
				return next(ctx, w, r)
			}
		}

		req, err := http.NewRequest("GET", "/middleware", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		api.Handle("GET", "/middleware", mockHandler, testMiddleware)
		api.ServeHTTP(rr, req)

		if !middlewareCalled {
			t.Error("Middleware was not called")
		}
	})
}

func TestAPI_ServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	shutdown := make(chan os.Signal, 1)

	api := New(shutdown)
	api.Handle(http.MethodGet, "/", mockHandler) // Register the mock handler

	// Call ServeHTTP directly on the API instance
	api.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
