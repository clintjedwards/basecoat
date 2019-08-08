package metrics

import (
	"log"
	"net/http"

	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitPrometheusService starts a long running http prometheus endpoint
func InitPrometheusService(config *config.Config) {

	utils.StructuredLog(utils.LogLevelInfo, "starting metrics http service",
		map[string]string{"url": config.Metrics.Endpoint})

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(config.Metrics.Endpoint, nil))
}
