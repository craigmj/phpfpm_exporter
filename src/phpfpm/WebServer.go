package phpfpm

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func WebServer(listen string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listen, nil)
}
