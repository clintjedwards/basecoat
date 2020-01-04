package metrics

import (
	"net/http"

	"github.com/clintjedwards/basecoat/config"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitPrometheusService starts a long running http prometheus endpoint
func InitPrometheusService(config *config.Config) {

	zap.S().Infow("starting metrics http service",
		"url", config.Metrics.Endpoint)

	http.Handle("/metrics", promhttp.Handler())
	zap.S().Fatal(http.ListenAndServe(config.Metrics.Endpoint, nil))
}
