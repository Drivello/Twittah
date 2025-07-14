package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics for TweetService
var (
	TweetsCreationConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_creation_consumed_total",
			Help: "Total of tweets creation consumed",
		})

	TweetsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_created_total",
			Help: "Total of tweets created",
		})

	TweetsCreatedErrorTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_created_error_total",
			Help: "Total of tweets created with error",
		})

	TweetsDeletionConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_deletion_consumed_total",
			Help: "Total of tweets deletion consumed",
		})

	TweetsDeletedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_deleted_total",
			Help: "Total of tweets deleted",
		})

	TweetsDeletedErrorTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_deleted_error_total",
			Help: "Total of tweets deleted with error",
		})

	TimelineCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_timeline_cache_hits_total",
			Help: "Total of hits in the timeline cache",
		})

	TimelineCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tweets_timeline_cache_misses_total",
			Help: "Total of misses in the timeline cache",
		})
)

// Init registers all metrics for TweetService
func Init() {
	prometheus.MustRegister(TweetsCreationConsumedTotal, TweetsCreatedTotal, TweetsCreatedErrorTotal, TweetsDeletionConsumedTotal, TweetsDeletedTotal, TweetsDeletedErrorTotal, TimelineCacheHits, TimelineCacheMisses)
}

// Handler exposes the /metrics endpoint
func Handler() http.Handler {
	return promhttp.Handler()
}
