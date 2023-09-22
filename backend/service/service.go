package service

import (
	"errors"
	"ip2location-pfsense/cache"
	"ip2location-pfsense/config"
	"ip2location-pfsense/pfsense"
	"ip2location-pfsense/util"
	"ip2location-pfsense/webserve"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var bind_host string
var bind_port string
var ssl_cert string
var ssl_key string
var use_ssl bool
var ingest_logs string
var ip_requests string
var ip2l_results string
var ip2geomap string
var healthcheck string

// var valid_api_keys map[string]string

func init() {
	config.LoadConfigProvider("IP2Location-pfSense")
	bind_host = config.GetConfig().GetString("service.bind_host")
	bind_port = config.GetConfig().GetString("service.bind_port")
	use_ssl = config.GetConfig().GetBool("service.use_ssl")
	ssl_cert = config.GetConfig().GetString("service.ssl_cert")
	ssl_key = config.GetConfig().GetString("service.ssl_key")

	util.LogDebug("Initialising service and binding on %v:%v", bind_host, bind_port)
	ingest_logs = config.GetConfig().GetString("service.ingest_logs")
	util.LogDebug("Ingest logs: %v", ingest_logs)
	ip_requests = config.GetConfig().GetString("service.ip_requests")
	util.LogDebug("IP requests: %v", ip_requests)
	ip2l_results = config.GetConfig().GetString("service.ip2l_results")
	util.LogDebug("IP2Location results: %v", ip2l_results)
	ip2geomap = config.GetConfig().GetString("service.ip2geomap")
	util.LogDebug("IP2Location GeoMap: %v", ip2geomap)
	healthcheck = config.GetConfig().GetString("service.healthcheck")
	util.LogDebug("Health check: %v", healthcheck)
}

// Service is the main entry point for the service
// It starts the service and listens for requests
// It also handles the requests
func Start(args []string) {
	log.Print("Starting service ...")

	e := echo.New()
	e.HideBanner = true

	//e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "${time_rfc3339} ${id} ${remote_ip} ${method} ${uri} ${user_agent} ${status} ${error} ${latency} ${latency_human} ${bytes_in} ${bytes_out}\n"}))
	e.Logger.SetHeader("${time_rfc3339} ${id} ${remote_ip} ${method} ${uri} ${user_agent} ${status} ${error} ${latency} ${latency_human} ${bytes_in} ${bytes_out}\n")

	e.Use(middleware.Recover())
	e.GET(healthcheck, healthCheck)
	e.POST(ingest_logs, ingestLog)
	e.GET(ip2l_results, ip2Results)
	e.GET(ip2geomap, ip2MapResults)
	e.POST(ip_requests, ipRequest)
	e.GET(ip2geomap, ip2GeoMap)

	g := e.Group("/api")

	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			if key == "" {
				util.LogDebug("[service] Missing API key")
				return false, errors.New("missing api key")
			}
			valid_api_keys := config.GetConfig().GetStringSlice("apikeys")
			for _, valid_key := range valid_api_keys {
				if key == valid_key {
					util.LogDebug("[service] Valid API key recieved.")
					return true, nil
				}
			}
			return false, nil
		},
		ContinueOnIgnoredError: false,
	}))

	e = webserve.ServeEmbeddedErrorFiles(e)
	e = webserve.ServeErrorTemplate(e)
	e = webserve.ServeEmbedded(e)
	e = webserve.ServeRenderTemplate(e)

	e.HTTPErrorHandler = webserve.CustomHTTPErrorHandler

	util.LogDebug("[service] Service called with: %s", strings.Join(args, " "))

	useCache := config.GetConfig().GetBool("use_cache")
	if useCache {
		log.Print("Using Redis cache")
		cache.CreateInstances()
	}
	var err error

	log.Printf("[service] Binding to: %v port %v; using SSL: %v", bind_host, bind_port, use_ssl)
	if use_ssl {
		err = e.StartTLS(bind_host+":"+bind_port, ssl_cert, ssl_key)
	} else {
		err = e.Start(bind_host + ":" + bind_port)
	}

	util.HandleFatalError(err, "Failed to start service")
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
	filterLog := new(pfsense.FilterLog)
	if err := c.Bind(filterLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	util.LogDebug("[service] Received log entries\n")
	resultid := pfsense.ProcessLogEntries(*filterLog)

	return c.JSON(http.StatusOK, resultid)
}

// Returns IP2Location results
// Expected input is a resultid from a previous request, supplied as a query parameter
// e.g. http://localhost:9999/ip2lresults?id=xxxxxxxx
func ip2Results(c echo.Context) error {
	resultid := c.QueryParam("id")
	util.LogDebug("[service] Received request for resultid: %s\n", resultid)
	resultset, err := pfsense.GetResult(resultid)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, resultset)
}

// Static file handler for the IP2Location GeoMap
// Returns the HTML file for the GeoMap page
func ip2GeoMap(c echo.Context) error {
	return c.File("index.html")
}

// Responds to requests from the static geomap page
func ip2MapResults(c echo.Context) error {
	resultid := c.QueryParam("resultid")
	util.LogDebug("[service] Received request for resultid: %s\n", resultid)
	resultset, err := pfsense.GetRawResult(resultid)
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
	util.LogDebug("[service] Received request for IP2Location data\n")
	if err := c.Bind(pfLog); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, pfLog)
}
