package web

import (
	"fmt"
	"strings"

	"github.com/jpmchia/ip2location-pfsense/ip2location"
	"github.com/jpmchia/ip2location-pfsense/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/util"
	"github.com/labstack/echo/v4"
)

func Ip2lHandler(c echo.Context) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)
	var ip2lp *ip2location.Ip2LocationEntry
	var err error

	ip2lp, err = ip2location.RetrieveIpPlus(c.QueryParam("ip"))
	util.HandleError(err, "[web] Failed to retrieve IP2Location data for IP: %s", c.QueryParam("ip"))

	data["ip2l"] = ip2lp.IP

	data["Title"] = fmt.Sprintf("%s geolocation", ip2lp.IP)

	data["LocationTable"], err = constructLocationTable(ip2lp)
	util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))

	data["TechnicalTable"], err = constructTechnicalTable(ip2lp)
	util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))

	data["MetricsTable"], err = constructMetricsTable(ip2lp)
	util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))

	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["APIKey"] = generatedKey.Key
	data["Theme"] = c.QueryParam("theme")
	data["APIKeyExpires"] = generatedKey.Expires
	data["Lat"] = ip2lp.Latitude
	data["Lon"] = ip2lp.Longitude

	data["LocationClass"] = "active"
	data["WatchlistClass"] = "inactive"
	data["CacheClass"] = "inactive"
	data["ExportClass"] = "inactive"
	data["SettingsClass"] = "inactive"
	data["HelpClass"] = "inactive"

	data = IncludeShaders(data)

	return data, nil
}

func constructLocationTable(ip2lplus *ip2location.Ip2LocationEntry) (tableStr string, err error) {

	if ip2lplus == nil {
		return "", fmt.Errorf("ip2lplus is nil")
	}

	tableStr = "<table class=\"ip2l-table\">"

	if len(ip2lplus.AddressType) > 0 {
		tableStr += "<tr><th>Type</th><td>" + ip2lplus.AddressType + "</td></tr>"
	}

	if len(ip2lplus.CityName) > 0 {
		tableStr += "<tr><th>City</th><td>" + ip2lplus.CityName + "</td></tr>"
	}

	if len(ip2lplus.ZipCode) > 0 && ip2lplus.ZipCode != "-" {
		tableStr += "<tr><th>Zip Code</th><td>" + ip2lplus.ZipCode + "</td></tr>"
	}

	if len(ip2lplus.RegionName) > 0 {
		tableStr += "<tr><th>Region</th><td>" + ip2lplus.RegionName + "</td></tr>"
	}

	if len(ip2lplus.CountryName) > 0 {
		tableStr += "<tr><th>Country</th><td>" + ip2lplus.CountryName
	}

	if len(ip2lplus.CountryCode) > 0 {
		tableStr += "  (" + ip2lplus.CountryCode + ")</td></tr>"
	} else {
		tableStr += "</td></tr>"
	}

	if len(ip2lplus.Continent.Name) > 0 {
		tableStr += "<tr><th>Continent</th><td>" + ip2lplus.Continent.Name + "</td></tr>"
	}

	if len(ip2lplus.Continent.Hemisphere) > 0 {
		tableStr += fmt.Sprintf("<tr><th>Hemisphere</th><td>%s</td></tr>", strings.Join(ip2lplus.Continent.Hemisphere, ","))
	}

	if ip2lplus.Latitude != 0 {
		tableStr += fmt.Sprintf("<tr><th>Latitude</th><td>%E</td></tr>", ip2lplus.Latitude)
	}

	if ip2lplus.Longitude != 0 {
		tableStr += fmt.Sprintf("<tr><th>Latitude</th><td>%E</td></tr>", ip2lplus.Longitude)
	}

	if len(ip2lplus.TimeZoneInfo.Olson) > 0 {
		tableStr += fmt.Sprintf("<tr><th>Timezone</th><td>%s %s (%d)</td></tr>", ip2lplus.TimeZoneInfo.Olson, ip2lplus.TimeZoneInfo.CurrentTime, ip2lplus.TimeZoneInfo.GmtOffset)
	}

	tableStr += "</table>"

	return tableStr, nil
}

func constructTechnicalTable(ip2lplus *ip2location.Ip2LocationEntry) (tableStr string, err error) {

	if ip2lplus == nil {
		return "", fmt.Errorf("ip2lplus is nil")
	}

	tableStr = "<table class=\"ip2l-table\">"

	if len(ip2lplus.Asn) > 0 {
		tableStr += "<tr><th>ASN</th><td>" + ip2lplus.Asn + "</td></tr>"
	}

	if len(ip2lplus.As) > 0 {
		tableStr += "<tr><th>AS</th><td>" + ip2lplus.As + "</td></tr>"
	}

	if len(ip2lplus.Domain) > 0 {
		tableStr += "<tr><th>Domain</th><td>" + ip2lplus.Domain + "</td></tr>"
	}

	if len(ip2lplus.NetSpeed) > 0 {
		tableStr += "<tr><th>Net. Speed</th><td>" + ip2lplus.NetSpeed + "</td></tr>"
	}

	if len(ip2lplus.UsageType) > 0 {
		tableStr += "<tr><th>Type</th><td>" + ip2lplus.UsageType + "</td></tr>"
	}

	tableStr += "</table>"

	return tableStr, nil
}

func constructMetricsTable(ip2lplus *ip2location.Ip2LocationEntry) (tableStr string, err error) {
	// var data = make(map[string]interface{})
	// var metrics *cache.Metrics

	tableStr = "<table class=\"ip2l-table\">"

	if len(ip2lplus.FirstSeen) > 0 {
		tableStr += "<tr><th>First seen</th><td>" + ip2lplus.FirstSeen + "</td></tr>"
	}

	if len(ip2lplus.LastSeen) > 0 {
		tableStr += "<tr><th>Last seen</th><td>" + ip2lplus.LastSeen + "</td></tr>"
	}

	if ip2lplus.Hits > 0 {
		tableStr += fmt.Sprintf("<tr><th>No. of hits</th><td>%d</td></tr>", ip2lplus.Hits)
	}

	tableStr += "</table>"

	return tableStr, nil
}
