package metrics

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitPrometheusService starts a long running http prometheus endpoint
func InitPrometheusService(endpoint string) {
	zap.S().Infow("starting metrics http service",
		"url", endpoint)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(endpoint, nil))
}
