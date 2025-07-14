package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// UserService metrics
var (
	FollowsConsumerTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "follows_consumer_total",
			Help: "Total of follows consumed by the service.",
		},
	)

	FollowsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "follows_created_total",
			Help: "Total of follows created by the service.",
		},
	)

	FollowsErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "follows_errors_total",
			Help: "Total of errors in follow operations.",
		},
	)

	UnfollowsConsumerTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "unfollows_consumer_total",
			Help: "Total of unfollows consumed by the service.",
		},
	)

	UnfollowsDoneTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "unfollows_done_total",
			Help: "Total of unfollows done by the service.",
		},
	)

	UnfollowsErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "unfollows_errors_total",
			Help: "Total of errors in unfollow operations.",
		},
	)
)

// Init registers UserService metrics
func Init() {
	prometheus.MustRegister(
		FollowsCreatedTotal,
		UnfollowsConsumerTotal,
		UnfollowsDoneTotal,
		UnfollowsErrorsTotal,
	)
}

// Handler returns the Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
