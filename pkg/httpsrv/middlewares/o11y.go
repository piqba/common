package middlewares

import (
	"github.com/piqba/common/pkg/o11y"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MetricCounter ...
func MetricCounter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteAddrToParse, ip := "", ""
			switch {
			case strings.Contains(r.RemoteAddr, "[::1]"):
				remoteAddrToParse = strings.Replace(r.RemoteAddr, "[::1]", "localhost", -1)
				ip = strings.Split(remoteAddrToParse, ":")[0]
			default:
				ip = strings.Split(r.RemoteAddr, ":")[0]
			}

			o11y.MetricCounterByEndpoint.With(prometheus.Labels{
				"path":   r.RequestURI,
				"method": r.Method,
				"ip":     ip,
			}).Inc()
			next.ServeHTTP(w, r) // dispatch the request

		})
	}
}

// MetricHistogram ...
func MetricHistogram() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				o11y.MetricHttp.
					WithLabelValues(
						r.URL.Path,
						r.Method,
						strconv.Itoa(ww.Status()),
						strconv.Itoa(ww.BytesWritten()),
					).
					Observe(time.Since(start).Seconds())
			}()
			next.ServeHTTP(ww, r) // dispatch the request

		})
	}
}
