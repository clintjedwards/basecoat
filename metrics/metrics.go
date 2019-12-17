package metrics

import (
	"log"
	"net/http"

	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/toolkit/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitPrometheusService starts a long running http prometheus endpoint
func InitPrometheusService(config *config.Config) {

	logger.Log().Infow("starting metrics http service",
		"url", config.Metrics.Endpoint)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(config.Metrics.Endpoint, nil))
}
