package rest

import (
	"context"
	"net/http"
	"os"
	"syscall"
)

// A Handler is a type that handles a http request within the framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// API is the handler for api package.
type API struct {
	shutdown chan os.Signal
	mux      *http.ServeMux
	mw       []Middleware
}

// New creates an API struct with provided middleware.
//
//	shutdown: channel to signal shutdown
//	mw: list of middleware to execute on each request
func New(shutdown chan os.Signal, mw ...Middleware) *API {
	return &API{
		shutdown: shutdown,
		mux:      http.NewServeMux(),
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *API) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *API) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// First wrap handler specific middleware around this handler
	handler = wrapMiddleware(mw, handler)

	// Add the package's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the context with the required values to
		// process the request.
		ctx := context.WithValue(r.Context(), key, &ContextValues{})

		// Register this path
		_ = SetPath(ctx, path)

		// Execute the handler.
		//
		// If there is an error, then shutdown the server.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
		}
	})

	a.mux.Handle(method+" "+path, h1)
}

// ServeHTTP implements the http.Handler interface. It's the entry point for
// all http traffic.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
