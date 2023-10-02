package web

import (
	"fmt"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/service/apikey"
	"github.com/labstack/echo/v4"
)

func WatchListHandler(c echo.Context) error {
	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)

	data["Title"] = "IP2Location.io Backend service for pfSense"
	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["APIKey"] = generatedKey.Key
	data["Theme"] = c.QueryParam("theme")
	data["APIKeyExpires"] = generatedKey.Expires
	data["Message"] = "Hello there!"
	data = IncludeShaders(data)

	fmt.Printf("ContentHandler: Rendering template with: %s", data)
	return c.Render(http.StatusOK, "watchlist.html.tmpl", data)
}

func constructLogHistoryTable(ipAddr string) (tableStr string, err error) {

	// var logEntries []pfsense.Ip2Map = pfsense.GetLogEntries(ipAddr)

	tableStr = "<table class=\"ip2l-table\">"

	tableStr += "<tr><th>Time</th><th>Dir.</th><th>IF</th><th>Proto</th><th>Dest.</th><th>Port</th><th>IF</th><th>Act</th><th>Reason</th></tr>"
	// tableStr += "<tr><td>Time</td><td>Dir.</td><td>IF</td><td>Proto</td><td>Dest.</td><td>Port</td><td>IF</td><td>Act</td><td>Reason</td></tr>"
	// for _, entry := range logEntries {
	// 	tableStr += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", entry.Time, entry.Direction, entry.Interface, entry.Proto, entry.Dstip, entry.Dstport, entry.Act, entry.Reason)
	// }

	tableStr += "</table>"

	return tableStr, nil
}
