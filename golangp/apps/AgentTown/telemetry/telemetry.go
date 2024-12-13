package telemetry

import (
	"fmt"
	"time"
)

// Telemetry represents the telemetry data for monitoring
type Telemetry struct {
	metrics map[string]int
}

// NewTelemetry creates a new telemetry instance
func NewTelemetry() *Telemetry {
	return &Telemetry{
		metrics: make(map[string]int),
	}
}

// RecordMetric records a metric
func (t *Telemetry) RecordMetric(name string, value int) {
	t.metrics[name] = value
}

// ReportMetrics reports the collected metrics
func (t *Telemetry) ReportMetrics() {
	for name, value := range t.metrics {
		fmt.Printf("Metric %s: %d\n", name, value)
	}
}

// Monitor starts monitoring metrics
func (t *Telemetry) Monitor() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			t.ReportMetrics()
		}
	}()
}
