package routes

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/config"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/labstack/echo/v4"
)

var HealthCheck_Route string

func init() {

	util.LogDebug("[healthcheck] Initialising health check service ...")

	config.Configure()

	conf := config.GetConfiguration().Service
	HealthCheck_Route = conf.HealthCheck

	util.Log("[healthcheck] Health check endpoint: %v", HealthCheck_Route)
}

// Health Check API
// Returns a simple string to indicate that the service is available
func HealthCheck_Handler(c echo.Context) error {
	util.Log("[healthcheck] Responding to health check request\n")
	return c.String(http.StatusOK, "Service is available.")
}
