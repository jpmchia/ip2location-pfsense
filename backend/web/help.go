package web

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/ip2location"
	"github.com/jpmchia/ip2location-pfsense/backend/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

func HelpHandler(c echo.Context) error {
	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)
	var ip2lp *ip2location.Ip2LocationEntry
	var err error

	ip2lp, err = ip2location.RetrieveIpPlus(c.QueryParam("ip"))
	util.HandleError(err, "[web] Failed to retrieve IP2Location data for IP: %s", c.QueryParam("ip"))

	data["ip2l"] = ip2lp.IP
	data["Title"] = "About IP2Location-pfSense"
	data["HelpContent"] = ""
	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["APIKey"] = generatedKey.Key
	data["Theme"] = c.QueryParam("theme")
	data["APIKeyExpires"] = generatedKey.Expires
	data["Lat"] = ip2lp.Latitude
	data["Lon"] = ip2lp.Longitude

	data["HelpContent"] = "Help"

	data["LocationClass"] = "inactive"
	data["WatchlistClass"] = "inactive"
	data["CacheClass"] = "inactive"
	data["ExportClass"] = "inactive"
	data["SettingsClass"] = "inactive"
	data["HelpClass"] = "active"

	data = IncludeShaders(data)

	util.LogDebug("ContentHandler: Rendering template with: %s", data)

	return c.Render(http.StatusOK, "help.html.tmpl", data)
}
