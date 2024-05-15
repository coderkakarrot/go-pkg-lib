package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coderkakarrots/go-pkg-lib/metric"
)

const meterName = "meterName"

// Handler is a custom handler that will provide handlers.
type Handler struct {
	grMetric  *metric.Gauge
	dqlMetric *metric.Histogram
}

// HelloWorld is a simple handler that returns "Hello, world!".
func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	// Simulate query time between 0 and 0.9 seconds
	t := rand.Float64() * 0.9
	// Record database metric.
	h.dqlMetric.Record(r.Context(), t, nil)

	go func() {
		opt := metric.Attributes{
			"handler": "goroutine",
		}

		// Record goroutine metric by increment the count.
		h.grMetric.Add(r.Context(), 1, opt)
		// Decrement the goroutine count.
		defer h.grMetric.Add(r.Context(), -1, opt)

		time.Sleep(2 * time.Second)
	}()

	w.Write([]byte("Hello, world!"))
}

func main() {
	//-------------------------------------
	// Create a new HTTP server
	//-------------------------------------

	mux := http.NewServeMux()

	// Create a channel to listen for the interrupt signal.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	//-------------------------------------
	// Create a new metric provider
	//-------------------------------------
	metrics, err := metric.Initialise(meterName)
	if err != nil {
		log.Fatalf("failed to initialize provider: %v", err)
		return
	}

	//------------------------------------------
	// Create a new histogram with custom bounds
	//------------------------------------------
	// Define custom bucket boundaries.
	boundaries := []float64{0.1, 0.5, 1, 2, 5, 10}

	// Create a new histogram with the given boundaries (otherwise default will be used).
	dqlMetric, err := metrics.NewHistogram("db_query_latency", "Latency of database queries", boundaries...)
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
		return
	}

	//-------------------------------------
	// Register the handlers
	//-------------------------------------

	// Create a new handler.
	h := &Handler{
		grMetric:  metrics.Goroutine,
		dqlMetric: dqlMetric,
	}

	// Register the home page handler.
	mux.HandleFunc("/", h.HelloWorld)

	// Register the metric handler to see the result.
	mux.HandleFunc("/metrics", metrics.Handler().ServeHTTP)

	server := &http.Server{
		Addr:    ":8080",
		Handler: MetricMiddleware(metrics, mux),
	}

	log.Println("Server is running on port 8080")
	log.Println("Access http://localhost:8080 and see the metrics at http://localhost:8080/metrics")
	go server.ListenAndServe()

	<-sig
	log.Println("Shutting down the server...")

	// Create a deadline for the outstanding requests.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}

// MetricMiddleware is a middleware that records metrics for each request.
func MetricMiddleware(m *metric.Metric, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			opt := metric.Attributes{
				"method": r.Method,
				"path":   r.URL.Path,
			}

			m.Latency.Record(r.Context(), elapsed.Seconds(), opt)
		}()

		ctx := r.Context()
		opt := metric.Attributes{
			"method": r.Method,
			"path":   r.URL.Path,
		}

		m.Request.Add(ctx, 1, opt)

		handler.ServeHTTP(w, r)
	})
}
