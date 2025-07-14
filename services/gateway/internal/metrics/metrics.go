package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GatewayService metrics
var (
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requests HTTP recibidos por método y código de estado.",
		},
		[]string{"method", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duración de requests HTTP en segundos.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Cantidad de requests HTTP concurrentes.",
		},
	)

	KafkaProducedMessagesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_produced_messages_total",
			Help: "Total de mensajes producidos a Kafka.",
		},
	)
)

// Init registers GatewayService metrics
func Init() {
	prometheus.MustRegister(
		HTTPRequestTotal,
		HTTPRequestDuration,
		ActiveRequests,
		KafkaProducedMessagesTotal,
	)
}

// Handler returns the Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
