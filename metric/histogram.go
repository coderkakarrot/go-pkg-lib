package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// Histogram is a metric that samples observations (usually things like request durations or response sizes)
// and counts them in configurable buckets.
type Histogram struct {
	histogram metric.Float64Histogram
}

// NewHistogram creates a new histogram with the given name, description and bounds.
func (m *Metric) NewHistogram(name, description string, bounds ...float64) (*Histogram, error) {
	h, err := m.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(bounds...),
	)
	if err != nil {
		return nil, err
	}

	return &Histogram{
		histogram: h,
	}, nil
}

// Record adds the given value in the given histogram.
func (h *Histogram) Record(ctx context.Context, incr float64, attr Attributes) {
	h.histogram.Record(ctx, incr, metric.WithAttributes(attr.toOtel()...))
}
