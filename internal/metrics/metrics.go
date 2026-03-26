package metrics

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Metrics struct {
	mu                 sync.Mutex
	startedAt          time.Time
	totalRequests      map[string]int64
	totalErrors        int64
	totalRateLimited   int64
	currentInFlight    int64
	totalDurationNanos map[string]int64
}

func New() *Metrics {
	return &Metrics{
		startedAt:          time.Now(),
		totalRequests:      make(map[string]int64),
		totalDurationNanos: make(map[string]int64),
	}
}

func (m *Metrics) IncRequests(network string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests[network]++
}

func (m *Metrics) IncErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalErrors++
}

func (m *Metrics) IncRateLimited() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRateLimited++
}

func (m *Metrics) IncInFlight() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentInFlight++
}

func (m *Metrics) DecInFlight() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentInFlight > 0 {
		m.currentInFlight--
	}
}

func (m *Metrics) AddDuration(network string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalDurationNanos[network] += duration.Nanoseconds()
}

func (m *Metrics) ExportPrometheus() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var builder strings.Builder

	builder.WriteString("# HELP dns_requests_total Total DNS requests.\n")
	builder.WriteString("# TYPE dns_requests_total counter\n")
	for network, total := range m.totalRequests {
		builder.WriteString(fmt.Sprintf("dns_requests_total{network=%q} %d\n", network, total))
	}

	builder.WriteString("# HELP dns_request_errors_total Total DNS request errors.\n")
	builder.WriteString("# TYPE dns_request_errors_total counter\n")
	builder.WriteString(fmt.Sprintf("dns_request_errors_total %d\n", m.totalErrors))

	builder.WriteString("# HELP dns_rate_limited_total Total rate limited DNS requests.\n")
	builder.WriteString("# TYPE dns_rate_limited_total counter\n")
	builder.WriteString(fmt.Sprintf("dns_rate_limited_total %d\n", m.totalRateLimited))

	builder.WriteString("# HELP dns_inflight_requests Current in-flight DNS requests.\n")
	builder.WriteString("# TYPE dns_inflight_requests gauge\n")
	builder.WriteString(fmt.Sprintf("dns_inflight_requests %d\n", m.currentInFlight))

	builder.WriteString("# HELP dns_request_duration_seconds_total Total request duration in seconds.\n")
	builder.WriteString("# TYPE dns_request_duration_seconds_total counter\n")
	for network, total := range m.totalDurationNanos {
		builder.WriteString(fmt.Sprintf("dns_request_duration_seconds_total{network=%q} %.6f\n", network, float64(total)/float64(time.Second)))
	}

	builder.WriteString("# HELP dns_uptime_seconds DNS forwarder uptime in seconds.\n")
	builder.WriteString("# TYPE dns_uptime_seconds gauge\n")
	builder.WriteString(fmt.Sprintf("dns_uptime_seconds %.0f\n", time.Since(m.startedAt).Seconds()))

	return builder.String()
}
