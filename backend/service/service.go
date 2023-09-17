package service

import (
	"ip2location-pfsense/cache"
	"ip2location-pfsense/config"
	"ip2location-pfsense/pfsense"
	. "ip2location-pfsense/pfsense"
	. "ip2location-pfsense/util"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var bind_host string
var bind_port string

func init() {
	config.LoadConfigProvider("IP2Location-pfSense")
	bind_host = config.GetConfig().GetString("service.bind_host")
	bind_port = config.GetConfig().GetString("service.bind_port")
	LogDebug("Initialising service and binding on %s:%s", bind_host, bind_port)
}

// Service is the main entry point for the service
// It starts the service and listens for requests
// It also handles the requests
func Start(args []string) {
	log.Print("Starting service ...")

	e := echo.New()
	e.POST("/filterlog", ingestLog)
	e.POST("/ip2location", ipRequest)
	e.POST("/ip2results", ip2Results)
	e.GET("/ip2geomap", ip2MapResults)
	e.GET("/health", healthCheck)

	LogDebug("Service called with: %s\n", strings.Join(args, " "))

	// useCacheStr := strconv.FormatBool(config.GetConfig().GetBool("use_cache"))
	useCache := config.GetConfig().GetBool("use_cache")
	if useCache {
		log.Print("Using Redis cache")
		cache.CreateInstances()
	}

	log.Printf("Binding to: %v port %v; using cache: %v", bind_host, bind_port, useCache)

	err := e.Start(bind_host + ":" + bind_port)
	HandleFatalError(err, "Failed to start service")
}

// Health Check API
// Returns a simple string to indicate that the service is available
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Service is available.")
}

// Process pfSense Filter Logs
// Expected input is a pfSense FilterLog JSON object (see pfsense/pfsense.go)
// Returns a JSON object with the IP address and the IP2Location data
func ingestLog(c echo.Context) error {
	filterLog := new(FilterLog)
	if err := c.Bind(filterLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	LogDebug("Received log entries\n")
	resultid := pfsense.ProcessLogEntries(*filterLog)

	return c.JSON(http.StatusOK, resultid)
}

func ip2Results(c echo.Context) error {
	resultid := c.QueryParam("resultid")
	LogDebug("Received request for resultid: %s\n", resultid)
	resultset, err := pfsense.GetResult(resultid)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, resultset)
}

// Returns IP2Location results
// Expected input is a resultid from a previous request, supplied as a query parameter
// e.g. http://localhost:9999/ip2geomap?resultid=1
// Returns a JSON object with the IP address and the IP2Location data
func ip2MapResults(c echo.Context) error {
	resultid := c.QueryParam("resultid")
	LogDebug("Received request for resultid: %s\n", resultid)
	resultset, err := pfsense.GetResult(resultid)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, resultset)
}

// Responds to a request for IP2Location data
// Expected input is a JSON object with a single IP address
// Returns a JSON object with the IP address and the IP2Location data
func ipRequest(c echo.Context) error {
	pfLog := new(pfsense.FilterLog) // Bind
	LogDebug("Received request for IP2Location data\n")
	if err := c.Bind(pfLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, pfLog)
}
