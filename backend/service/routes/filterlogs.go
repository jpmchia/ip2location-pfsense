package routes

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/config"
	"github.com/jpmchia/ip2location-pfsense/pfsense"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/labstack/echo/v4"
)

var FilterLogs_PostRoute string
var FilterLogs_GetRoute string

func init() {
	util.LogDebug("[filterlogs] Initialising filter logs service ...")

	config.Configure()
	conf := config.GetConfiguration().Service
	FilterLogs_PostRoute = conf.IngestLogs
	FilterLogs_GetRoute = conf.Results
	util.Log("[filterlogs] Ingest filter logs endpoint: %s", FilterLogs_PostRoute)
	util.Log("[filterlogs] Results endpoint: %s", FilterLogs_GetRoute)
}

// Process pfSense Filter Logs
// Expected input is a pfSense FilterLog JSON object (see pfsense/pfsense.go)
// Returns a JSON object with the IP address and the IP2Location data
func PostLogsHandler(c echo.Context) error {
	filterLog := new(pfsense.FilterLog)
	if err := c.Bind(filterLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	util.Log("[filterlogs] Received log entries")
	resultid := pfsense.ProcessLogEntries(*filterLog)

	return c.JSON(http.StatusOK, resultid)
}

// Returns IP2Location results
// Expected input is a resultid from a previous request, supplied as a query parameter
// e.g. http://localhost:9999/ip2lresults?id=xxxxxxxx
func GetResultsHandler(c echo.Context) error {
	resultid := c.QueryParam("id")
	util.Log("[filterlogs] Received request for resultid: %s", resultid)
	resultset, err := pfsense.GetResult(resultid)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, resultset)
}
