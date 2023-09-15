package service

import (
	"log"
	"net/http"
	"pfSense/config"
	"pfSense/pfsense"
	. "pfSense/pfsense"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var bind_host string
var bind_port string

func init() {
	config.LoadConfigProvider("IP2LOCATION")

	log.Print("Configuring service ...")

	bind_host = config.Config().GetString("bind_host")
	bind_port = config.Config().GetString("bind_port")
}

// Service is the main entry point for the service
// It starts the service and listens for requests
// It also handles the requests
func Service(args []string) {
	log.Print("Starting service ...")

	e := echo.New()
	e.POST("/filterlog", ingestLog)
	e.POST("/ip2location", ipRequest)
	e.GET("/ip2geomap", ip2MapResults)
	e.GET("/health", healthCheck)

	log.Printf("Service called with: %s\n", strings.Join(args, " "))

	useCacheStr := strconv.FormatBool(config.Config().GetBool("use_cache"))
	log.Printf("Binding to " + bind_host + ":" + bind_port + "; Using cache: " + useCacheStr)

	e.Logger.Fatal(e.Start(bind_host + ":" + bind_port))
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

	resultid := pfsense.ProcessLogEntries(*filterLog)

	return c.JSON(http.StatusOK, resultid)
}

// Returns IP2Location results
// Expected input is a resultid from a previous request, supplied as a query parameter
// e.g. http://localhost:9999/ip2geomap?resultid=1
// Returns a JSON object with the IP address and the IP2Location data
func ip2MapResults(c echo.Context) error {
	resultid := c.QueryParam("resultid")
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
	if err := c.Bind(pfLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, pfLog)
}
