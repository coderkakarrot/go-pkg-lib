//go:build unit
// +build unit

package metric

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestMetric(t *testing.T) {
	t.Run("AddRequest", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		opt := Attributes{
			"A": "B",
			"C": "D",
		}
		ctx := context.Background()
		m.Request.Add(ctx, 1, opt)
	})

	t.Run("CustomCounter", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		ctx := context.Background()
		newC, _ := m.NewCounter("concurrent_request", "concurrent request count")

		newC.Add(ctx, 1, nil)

		// Add 2 to the counter.
		newC.Add(ctx, 2, nil)
	})

	t.Run("RecordLatency", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		ctx := context.Background()
		start := time.Now()
		elapsed := time.Since(start).Seconds()
		m.Latency.Record(ctx, elapsed, nil)
	})

	t.Run("CustomHistogram", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		ctx := context.Background()
		// Define custom bucket boundaries.
		boundaries := []float64{0.1, 0.5, 1, 2, 5, 10}

		newH, _ := m.NewHistogram("database_query_duration", "database query duration", boundaries...)

		// Simulate query time between 0 and 0.9 seconds
		queryTime := rand.Float64() * 0.9
		newH.Record(ctx, queryTime, nil)
	})

	t.Run("RecordGoroutine", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		ctx := context.Background()
		m.Goroutine.Add(ctx, 1, Attributes{
			"handler": "TestRecordGoroutine",
		})

		m.Goroutine.Add(ctx, -1, Attributes{
			"handler": "TestRecordGoroutine",
		})
	})

	t.Run("CustomGauge", func(t *testing.T) {
		t.Log("Initiated Metric Handler")
		m, err := Initialise("testMeter")
		if err != nil {
			t.Errorf("failed to initialize provider: %v", err)
		}

		ctx := context.Background()

		newG, _ := m.NewGauge("memory_usage", "memory usage gauge")
		// Increase gauge by 100.
		newG.Add(ctx, 100, nil)

		// Decrease gauge by 50.
		newG.Add(ctx, -50, nil)
	})
}
