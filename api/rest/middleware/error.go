package middleware

import (
	// Standard library packages
	"context"
	"fmt"
	"net/http"
	// Pantheon internal package
	"github.com/pantheon-systems/go-pkg-lib/api/rest"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
func Errors() rest.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler rest.Handler) rest.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := rest.Handler(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Run the next handler and catch any propagated error.
			err := handler(ctx, w, r)
			if err != nil {
				_ = rest.SetStatusCode(ctx, http.StatusInternalServerError)
				_ = rest.SetIsError(ctx)
				// Respond with the error back to the client
				if err := rest.Respond(ctx, w, err.Error(), http.StatusInternalServerError); err != nil {
					return fmt.Errorf("error responding with error: %w", err)
				}
			}

			return nil
		})

		return h
	}

	return m
}
