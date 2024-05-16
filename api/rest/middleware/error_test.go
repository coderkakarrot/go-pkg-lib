//go:build unit
// +build unit

package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	// Pantheon internal package
	"github.com/coderkakarrot/go-pkg-lib/api/rest"
)

func TestErrorsMiddleware(t *testing.T) {

	// Create a mock handler that returns an error
	mockHandlerWithError := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("simulated handler error")
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	shutdown := make(chan os.Signal, 1)

	api := rest.New(shutdown)

	api.Handle(http.MethodGet, "/", mockHandlerWithError, Errors()) // Register the mock handler

	// Call ServeHTTP directly on the API instance
	api.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusInternalServerError)
	}
	body := rr.Body.String()
	response := rest.Response{}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		t.Fatal("could not unmarshal response body")
	}
	if response.Success != false || response.Data != nil || response.Errors != "simulated handler error" {
		t.Errorf("handler returned wrong response: got %v", body)
	}
}
