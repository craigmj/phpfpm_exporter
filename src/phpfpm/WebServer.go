package phpfpm

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// WebServer listens on the provided address and handles the metrics endpoint
// for prometheus
func WebServer(listen string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listen, nil)
}
