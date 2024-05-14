//go:build unit
// +build unit

package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	t.Run("Success Response", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.WithValue(context.Background(), key, &ContextValues{})
		data := map[string]string{"message": "success"}
		err := Respond(ctx, rr, data, http.StatusOK)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check the status code and content type
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, status)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		// Decode and check the response body
		var resp Response
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if !resp.Success || resp.Timestamp == 0 || resp.Errors != nil {
			t.Errorf("Unexpected response content: %+v", resp)
		}
	})

	t.Run("Error Response", func(t *testing.T) {
		rr := httptest.NewRecorder()
		ctx := context.WithValue(context.Background(), key, &ContextValues{IsError: true})
		errData := "some error"

		err := Respond(ctx, rr, errData, http.StatusInternalServerError)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check the status code and content type
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Expected status %v, got %v", http.StatusInternalServerError, status)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		// Decode and check the response body
		var resp Response
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if resp.Success || resp.Timestamp == 0 || resp.Data != nil || resp.Errors != errData {
			t.Errorf("Unexpected response content: %+v", resp)
		}
	})
}
