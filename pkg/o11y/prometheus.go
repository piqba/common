package o11y

import (
	"fmt"
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
