package metric

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Metric is a wrapper around the OpenTelemetry metric provider.
type Metric struct {
	// provider is the OpenTelemetry metric provider.
	provider *sdkmetric.MeterProvider
	// meter is the OpenTelemetry meter.
	meter metric.Meter

	// Request is to store the request count.
	Request *Counter
	// Latency is to store the response time.
	Latency *Histogram
	// Goroutine is to store the concurrent goroutine count.
	Goroutine *Gauge
	// Errors is to store the error count.
	Errors *Counter
	// Panics is to store the panic count.
	Panics *Counter
}

// Initialise creates a new meter provider with the given meter name.
func Initialise(meterName string) (*Metric, error) {
	// Prometheus is the default exporter at the moment for this package.
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// Register the exporter.
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	// Get the meter.
	meter := provider.Meter(meterName)

	// Create a new metric struct.
	m := &Metric{
		provider: provider,
		meter:    meter,
	}

	// Create the counters.
	m.Request, err = m.NewCounter("request", "Incremental counter of all requests")
	if err != nil {
		return nil, err
	}

	m.Latency, err = m.NewHistogram("latency", "Measurement of request latencies")
	if err != nil {
		return nil, err
	}

	m.Goroutine, err = m.NewGauge("goroutine", "Gauge of all goroutines")
	if err != nil {
		return nil, err
	}

	m.Errors, err = m.NewCounter("error", "Incremental counter of all errors")
	if err != nil {
		return nil, err
	}

	m.Panics, err = m.NewCounter("panic", "Incremental counter of all panics")
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Shutdown shuts down the metric provider.
func (m *Metric) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}

// Delete removes the topic from pubsub
// using the context provided
func delete(ctx context.Context) error {
	return fmt.Errorf("topic '%s' does not exist", uuid.New().String())
}
