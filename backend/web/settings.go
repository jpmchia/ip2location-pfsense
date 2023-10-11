package web

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/ip2location"
	"github.com/jpmchia/ip2location-pfsense/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/labstack/echo/v4"
)

func SettingsHandler(c echo.Context) error {
	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)
	var ip2lp *ip2location.Ip2LocationEntry
	var err error

	ip2lp, err = ip2location.RetrieveIpPlus(c.QueryParam("ip"))
	util.HandleError(err, "[web] Failed to retrieve IP2Location data for IP: %s", c.QueryParam("ip"))

	data["ip2l"] = ip2lp.IP

	data["Title"] = "Settings and configuration"

	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["APIKey"] = generatedKey.Key
	data["Theme"] = c.QueryParam("theme")
	data["APIKeyExpires"] = generatedKey.Expires
	data["Lat"] = ip2lp.Latitude
	data["Lon"] = ip2lp.Longitude

	data["LocationClass"] = "inactive"
	data["WatchlistClass"] = "inactive"
	data["CacheClass"] = "inactive"
	data["ExportClass"] = "inactive"
	data["SettingsClass"] = "active"
	data["HelpClass"] = "inactive"

	data = IncludeShaders(data)

	util.LogDebug("ContentHandler: Rendering template with: %s", data)

	return c.Render(http.StatusOK, "settings.html.tmpl", data)
}
