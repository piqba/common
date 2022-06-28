package o11y

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// ExposeMetricServer expose default metrics from our golang process
func ExposeMetricServer(options PrometheusOptions) {
	http.Handle(options.RoutePath, promhttp.Handler())
	port := fmt.Sprintf(":%d", options.Port)
	log.Fatal(http.ListenAndServe(port, nil).Error())
}

var (
	MetricCounterByEndpoint = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kraken_app",
			Name:      "count_request_path",
			Help:      "Counts request by path",
		},
		[]string{
			"path",
			"method",
			"ip",
		},
	)
	MetricHttp = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"path", "method", "code", "bytes"})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(MetricCounterByEndpoint)
	err := prometheus.Register(MetricHttp)
	if err != nil {
		log.Fatal(err)
	}
}
