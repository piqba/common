# O11y module

This module, pkg or folder contain customs metrics and expose all logic for expose a prometheus instance for collect metrics from our applications.


```go
func ExposeMetricServer(options PrometheusOptions) {
	http.Handle(options.RoutePath, promhttp.Handler())
	port := fmt.Sprintf(":%d", options.Port)
	log.Fatal(http.ListenAndServe(port, nil).Error())
}

var (
	MetricByEndpoint = prometheus.NewCounterVec(
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
)
```