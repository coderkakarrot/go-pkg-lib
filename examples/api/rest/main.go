package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/pantheon-systems/go-pkg-lib/api/rest"
	"github.com/pantheon-systems/go-pkg-lib/api/rest/middleware"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	mw := make([]rest.Middleware, 0, 1)
	mw = append(mw, middleware.Errors())
	api := rest.New(shutdown, mw...)

	api.Handle("GET", "/hello", helloHandler)
	api.Handle("GET", "/error-handler-simulation", errorHandler)

	// Start the server (you might have a separate `Run` method in your framework)
	port := ":8080"
	fmt.Printf("Server listening on %s\n", port)
	if err := http.ListenAndServe(port, api); err != nil {
		api.SignalShutdown() // Signal shutdown in case of errors
	}
}

func helloHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := map[string]string{"message": "Hello from the Go REST API!"}
	err := rest.Respond(ctx, w, data, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.New("simulated handler error")
}
