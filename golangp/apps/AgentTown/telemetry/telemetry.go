package telemetry

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Telemetry represents the telemetry data for monitoring
type Telemetry struct {
	metrics map[string]*int64
}

var tm *Telemetry
var once sync.Once

// GetInstance returns the singleton instance of Telemetry
func GetInstance() *Telemetry {
	once.Do(func() {
		tm = &Telemetry{
			metrics: make(map[string]*int64),
		}
	})
	return tm
}

// RecordMetricInc records a metric atomically by adding the value to the current value
func RecordMetricInc(name string, value int64) {
	t := GetInstance()
	if _, exists := t.metrics[name]; !exists {
		var initialValue int64
		t.metrics[name] = &initialValue
	}
	atomic.AddInt64(t.metrics[name], value)
}

// RecordMetricSet records a metric atomically by setting the value
func RecordMetricSet(name string, value int64) {
	t := GetInstance()
	if _, exists := t.metrics[name]; !exists {
		var initialValue int64
		t.metrics[name] = &initialValue
	}
	atomic.StoreInt64(t.metrics[name], value)
}

// GetMetric retrieves the value of a metric atomically
func GetMetric(name string) int64 {
	t := GetInstance()
	if value, exists := t.metrics[name]; exists {
		return atomic.LoadInt64(value)
	}
	return 0
}

// ReportMetrics reports the collected metrics
func ReportMetrics() {
	t := GetInstance()
	for name, value := range t.metrics {
		fmt.Printf("[Telemetry] %s: %d\n", name, atomic.LoadInt64(value))
	}
}

// Monitor starts monitoring metrics
func Monitor(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			ReportMetrics()
		}
	}()
}
