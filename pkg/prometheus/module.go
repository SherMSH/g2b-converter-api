package prometheus

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.100, .250, .500, 1, 2.500, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	ServiceStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_status",
			Help: "Status of the service: 1 for up, 0 for down",
		},
	)

	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestProcessingTime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_processing_seconds",
			Help:       "Time spent processing requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "endpoint", "status"},
	)

	CPUUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage percentage",
		},
	)

	MemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	NumGoroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "num_goroutines",
			Help: "Number of active goroutines",
		},
	)

	Uptime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_uptime_seconds",
			Help: "Service uptime in seconds",
		},
	)

	SpecificErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "specific_error_total",
			Help: "Total occurrences of a specific error",
		},
	)
)

func Init() {
	prometheus.MustRegister(HttpRequestsTotal, RequestDuration, ServiceStatus, ErrorsTotal, RequestProcessingTime, CPUUsage, MemoryUsage, NumGoroutines, Uptime, SpecificErrorCounter)
}

func UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	CPUUsage.Set(float64(runtime.NumCPU()))
	MemoryUsage.Set(float64(memStats.Alloc))
	NumGoroutines.Set(float64(runtime.NumGoroutine()))
	SpecificErrorCounter.Inc()
	Uptime.Inc()
}
