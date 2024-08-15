package prometheus

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-ashland/internal/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func StartServer(promConfig *models.PrometheusConfig, logger models.Logger) {
	// Expose the metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	logger.Log(models.InfoLevel, "Starting the prometheus server", fmt.Sprintf("%s:%s", promConfig.Host, promConfig.Port))
	// Start the HTTP server
	err := http.ListenAndServe(fmt.Sprintf(":%s", promConfig.Port), nil)
	if err != nil {
		logger.Log(models.ErrorLevel, "Unable to start prometheus server", err)
		return
	}
}
