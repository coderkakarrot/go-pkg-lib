package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Response is the form used for API responses for success in the API.
type Response struct {
	// Success
	//
	Success bool `json:"success"`

	// Timestamp
	//
	// example: 1234567
	Timestamp int64 `json:"timestamp"`

	// Data
	// in: body
	Data interface{} `json:"data,omitempty"`

	// Errors
	// in: body
	Errors interface{} `json:"errors,omitempty"`
}

// Respond constructs and sends an HTTP response to the client.
// It handles both successful responses with data and error responses.
//
// Parameters:
//
//	ctx: The request context, containing values like error status and logging data.
//	w: The HTTP response writer used to send the response to the client.
//	data: The data to be included in the response (either data or errors).
//	statusCode: The HTTP status code to indicate the response type.
//
// Returns:
//
//	error: An error if something goes wrong during response preparation or sending.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	v, err := GetContextValues(ctx)
	if err != nil {
		return ErrMissingContext
	}

	// Set the status code for the request logger middleware in the context.
	err = SetStatusCode(ctx, statusCode)
	if err != nil {
		return err
	}

	// Create the response object to be sent back to the client.
	r := Response{
		Success:   !v.IsError,
		Timestamp: time.Now().UTC().Unix(),
	}

	// Check if the data is an error or not
	// If it is not an error, then set the data
	// If it is an error, then set the error
	if !v.IsError {
		r.Data = data
	} else {
		r.Errors = data
	}

	// Convert the response to json
	jd, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal fail: %w", err)
	}

	// set the content type now that we know there was no marshal error
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Send the result back to the client
	if _, err := w.Write(jd); err != nil {
		return fmt.Errorf("write fail: %w", err)
	}

	return nil
}
