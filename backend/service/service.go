package service

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/cache"
	"github.com/jpmchia/ip2location-pfsense/config"
	"github.com/jpmchia/ip2location-pfsense/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/service/routes"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/jpmchia/ip2location-pfsense/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebService struct {
	bind_host string
	bind_port string
	ssl_cert  string
	ssl_key   string
	UseSSL    bool
}

var webService WebService
var conf config.ServiceOptions

func init() {

	config.Configure()
	conf = config.GetConfiguration().Service

	webService.bind_host = conf.BindHost
	webService.bind_port = conf.BindPort
	webService.UseSSL = conf.UseSSL
	webService.ssl_cert = conf.SSLCert
	webService.ssl_key = conf.SSLKey

	util.Log("[service] Service host and port: %v:%v", webService.bind_host, webService.bind_port)
}

// Service is the main entry point for the service
// It starts the service and listens for requests
func Start(args []string) {

	util.Log("[service] Starting service ...")
	var err error
	// Create a new echo instance
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogProtocol:  true,
		LogRequestID: true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogRoutePath: true,
		LogURIPath:   true,
		LogStatus:    true,
		LogError:     true,
		BeforeNextFunc: func(c echo.Context) {
			//c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			util.Log("[service] REQUEST: %v, RemoteIP: %v, Host: %v, RoutePath: %v, Method: %v, URI: %v, Status: %v, Error: %v", v.Protocol, v.RemoteIP, v.Host, v.RoutePath, v.Method, v.URI, v.Status, v.Error)
			return nil
		},
	}))

	e.Logger.SetHeader("${time_rfc3339_nano} ${id} ${remote_ip} ${method} ${uri} ${user_agent} ${status} ${error} ${latency} ${latency_human} ${bytes_in} ${bytes_out}\n")

	// Recover from panics
	e.Use(middleware.Recover())

	// CORS
	hosts := conf.AllowHosts
	if conf.AllowHosts[0] != "*" {
		if len(conf.AllowHosts) > 1 {
			if conf.AllowHosts[0] != "*" {
				e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
					AllowOrigins: hosts,
					AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Authorization"},
					AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
				}))
			}
		}
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Authorization"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		}))
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Rewrite(map[string]string{
		"/": "/main.html",
	}))

	// Routes
	// Handler
	e.GET(routes.HealthCheck_Route, routes.HealthCheck_Handler) // Health Check
	e.GET("/api/key", apikey.ApiKeyHandler)                     // Get API key

	// pfSense Filter Logs endpoints
	e.POST(routes.FilterLogs_PostRoute, routes.PostLogsHandler) // Ingest pfSense Filter Logs
	e.GET(routes.FilterLogs_GetRoute, routes.GetResultsHandler) // Get IP2Location results

	// WatchList endpoints
	e.POST(routes.WatchList_PostItemRoute, routes.PostItemHandler)
	e.GET(routes.WatchList_GetItemRoute, routes.GetItemHandler)
	e.GET(routes.WatchList_GetRoute, routes.GetHandler)
	e.DELETE(routes.WatchList_DeleteItemRoute, routes.DeleteHandler)

	// Web content
	if conf.EnableWeb {
		util.Log("[service] Enabling web content: %s", conf.HomePage)
		e.GET("/ip2l", web.ContentHandler)
		e.GET(conf.HomePage, web.ContentHandler)
	} else {
		util.Log("[service] Web content is disabled. Use the --enable-web flag to enable it or set the enable_web option in the configuration file.")
	}

	e = web.ServeEmeddedContent(e)
	e = web.ServeEmbeddedTemplates(e)

	apiAuth := e.Group("/api")

	apiAuth.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:              "header:" + echo.HeaderAuthorization,
		AuthScheme:             "Bearer",
		Validator:              apikey.ValidateToken,
		ContinueOnIgnoredError: false,
	}))

	ipl2Auth := e.Group("/ip2l")

	ipl2Auth.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:              "query:key",
		Validator:              apikey.ValidateKey,
		ContinueOnIgnoredError: false,
	}))

	e = web.ServeFavIcons(e)

	e.HTTPErrorHandler = web.CustomHTTPErrorHandler

	useCache := config.GetConfiguration().UseRedis
	if useCache {
		util.LogDebug("[service] Using Redis cache")
		cache.CreateInstances()
	}

	util.Log("[service] Binding to: %v port %v; using SSL: %v", webService.bind_host, webService.bind_port, webService.UseSSL)

	routes := e.Routes()
	for _, route := range routes {
		util.Log("[service] Method: %s %v  ==>  %v", route.Method, route.Path, route.Name)
	}

	if webService.UseSSL {
		util.LogDebug("[service] Using SSL")
		util.LogDebug("[service] Certifcate: %v; Key: %v", webService.ssl_cert, webService.ssl_key)
		err = e.StartTLS(webService.bind_host+":"+webService.bind_port, webService.ssl_cert, webService.ssl_key)
	} else {
		err = e.Start(webService.bind_host + ":" + webService.bind_port)
	}

	util.HandleFatalError(err, "[service] Failed to start service")
}
