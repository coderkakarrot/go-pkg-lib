package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
type Gauge struct {
	gauge metric.Float64UpDownCounter
}

// NewGauge creates a new gauge with the given name and description.
func (m *Metric) NewGauge(name, description string) (*Gauge, error) {
	g, err := m.meter.Float64UpDownCounter(name, metric.WithDescription(description))
	if err != nil {
		return nil, err
	}

	return &Gauge{
		gauge: g,
	}, nil
}

// Add adds the given gauge to the given value.
func (g *Gauge) Add(ctx context.Context, value float64, attr Attributes) {
	g.gauge.Add(ctx, value, metric.WithAttributes(attr.toOtel()...))
}
