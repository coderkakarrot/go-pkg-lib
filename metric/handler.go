package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns the handler associated with the metrics object. Prometheus by default.
func (m *Metric) Handler() http.Handler {
	return promhttp.Handler()
}
