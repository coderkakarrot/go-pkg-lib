package metric

import "go.opentelemetry.io/otel/attribute"

// Attributes is a map of key-value pairs that can be used to add metadata to metrics.
type Attributes map[string]string

// toOtel converts the Attributes to OpenTelemetry attributes.
func (a Attributes) toOtel() []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(a))
	for k, v := range a {
		attrs = append(attrs, attribute.String(k, v))
	}

	return attrs
}
