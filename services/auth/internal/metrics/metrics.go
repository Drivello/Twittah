package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// AuthService metrics
var (
	AuthRegisterUserSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_register_user_success_total",
			Help: "Total de usuarios registrados exitosamente.",
		},
	)

	AuthRegisterUserFailureTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_register_user_failure_total",
			Help: "Total de usuarios registrados exitosamente.",
		},
	)

	DatabaseRegisterUserQueryDuration_seconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "auth_database_register_user_query_duration_seconds",
			Help: "Duration of database queries.",
		},
	)
)

// Init registers AuthService metrics
func Init() {
	prometheus.MustRegister(AuthRegisterUserSuccessTotal, AuthRegisterUserFailureTotal, DatabaseRegisterUserQueryDuration_seconds)
}

// Handler returns the Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
