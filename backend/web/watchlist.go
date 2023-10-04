package web

import (
	"github.com/jpmchia/ip2location-pfsense/backend/ip2location"
	"github.com/jpmchia/ip2location-pfsense/backend/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

func WatchListHandler(c echo.Context) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)
	var ip2lp *ip2location.Ip2LocationEntry
	var err error

	ip2lp, err = ip2location.RetrieveIpPlus(c.QueryParam("ip"))
	util.HandleError(err, "[web] Failed to retrieve IP2Location data for IP: %s", c.QueryParam("ip"))

	data["ip2l"] = ip2lp.IP
	data["Title"] = "Watch List"

	// data["LocationTable"], err = constructLocationTable(ip2lp)
	// util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))
	// data["TechnicalTable"], err = constructTechnicalTable(ip2lp)
	// util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))
	// data["MetricsTable"], err = constructMetricsTable(ip2lp)
	// util.HandleError(err, "[web] Failed to construct location table for IP: %s", c.QueryParam("ip"))

	data["CurrentWatch"], err = constructCurrentWatchListTable()
	util.HandleError(err, "constructCurrentWatchListTable")
	data["PreviousWatch"], err = constructHistoricWatchListTable()
	util.HandleError(err, "constructHistoricWatchListTable")

	util.Log("WatchListHandler:  CurrentWatch: %v", data)

	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["APIKey"] = generatedKey.Key
	data["Theme"] = c.QueryParam("theme")
	data["APIKeyExpires"] = generatedKey.Expires
	data["Lat"] = ip2lp.Latitude
	data["Lon"] = ip2lp.Longitude

	data["LocationClass"] = "inactive"
	data["WatchlistClass"] = "active"
	data["CacheClass"] = "inactive"
	data["ExportClass"] = "inactive"
	data["SettingsClass"] = "inactive"
	data["HelpClass"] = "inactive"

	data = IncludeShaders(data)

	return data, nil
}

func constructCurrentWatchListTable() (tableStr string, err error) {
	tableStr = "<table class=\"ip2l-table\">"
	tableStr += "<tr><th>Time</th><th>Dir.</th><th>IF</th><th>Proto</th><th>Dest.</th><th>Port</th><th>IF</th><th>Act</th><th>Reason</th></tr>"
	tableStr += "</table>"
	return tableStr, nil
}

func constructHistoricWatchListTable() (tableStr string, err error) {
	tableStr = "<table class=\"ip2l-table\">"
	tableStr += "<tr><th>Time</th><th>Dir.</th><th>IF</th><th>Proto</th><th>Dest.</th><th>Port</th><th>IF</th><th>Act</th><th>Reason</th></tr>"
	tableStr += "</table>"
	return tableStr, nil
}

// func constructLogHistoryTable(ipAddr string) (tableStr string, err error) {
// 	// var logEntries []pfsense.Ip2Map = pfsense.GetLogEntries(ipAddr)
// 	tableStr = "<table class=\"ip2l-table\">"
// 	tableStr += "<tr><th>Time</th><th>Dir.</th><th>IF</th><th>Proto</th><th>Dest.</th><th>Port</th><th>IF</th><th>Act</th><th>Reason</th></tr>"
// 	// tableStr += "<tr><td>Time</td><td>Dir.</td><td>IF</td><td>Proto</td><td>Dest.</td><td>Port</td><td>IF</td><td>Act</td><td>Reason</td></tr>"
// 	// for _, entry := range logEntries {
// 	// 	tableStr += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", entry.Time, entry.Direction, entry.Interface, entry.Proto, entry.Dstip, entry.Dstport, entry.Act, entry.Reason)
// 	// }
// 	tableStr += "</table>"
// 	return tableStr, nil
// }
