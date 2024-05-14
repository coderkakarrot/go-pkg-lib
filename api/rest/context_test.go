//go:build unit
// +build unit

package rest

import (
	"context"
	"net/http"
	"testing"
)

func TestGetContextValues(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), key, &ContextValues{
			StatusCode: http.StatusOK,
			IsError:    false,
			Path:       "/test",
		})

		v, err := GetContextValues(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if v.StatusCode != http.StatusOK || v.IsError != false || v.Path != "/test" {
			t.Errorf("Unexpected values: %+v", v)
		}
	})

	t.Run("Missing Context", func(t *testing.T) {
		_, err := GetContextValues(context.Background()) // Empty context
		if err != ErrMissingContext {
			t.Fatalf("Expected error %v, got %v", ErrMissingContext, err)
		}
	})
}

func TestSetStatusCode(t *testing.T) {
	ctx := context.WithValue(context.Background(), key, &ContextValues{})

	err := SetStatusCode(ctx, http.StatusNotFound)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	v, _ := GetContextValues(ctx)
	if v.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %v, got %v", http.StatusNotFound, v.StatusCode)
	}
}

func TestSetIsError(t *testing.T) {
	ctx := context.WithValue(context.Background(), key, &ContextValues{})

	err := SetIsError(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	v, _ := GetContextValues(ctx)
	if !v.IsError {
		t.Errorf("Expected IsError to be true")
	}
}

func TestSetPath(t *testing.T) {
	ctx := context.WithValue(context.Background(), key, &ContextValues{})

	err := SetPath(ctx, "/new-path")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	v, _ := GetContextValues(ctx)
	if v.Path != "/new-path" {
		t.Errorf("Expected path '/new-path', got '%s'", v.Path)
	}
}
