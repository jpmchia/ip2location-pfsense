package routes

import (
	"errors"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/ip2location"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

var Ip2Location_PostRoute = "/ip2l/ipaddr"
var Ip2Location_GetRoute = "/ip2l/ipaddr"

func init() {
	util.LogDebug("[ip2location] Initialising IP2Location service ...")
	util.Log("[ip2location] IP2Location post endpoint: %s", Ip2Location_PostRoute)
	util.Log("[ip2location] IP2Location get list endpoint: %s", Ip2Location_GetRoute)
}

func PostIp2LocationHandler(c echo.Context) error {
	ip_param := c.QueryParam("ip")
	key := c.QueryParam("key")
	util.Log("[ip2location] Received request for IP: %s\n", ip_param)

	err, ipAddr := ip2location.RetrieveIpLocationFromCache(ip_param, key)
	util.HandleError(err, "[ip2location] Failed to retrieve IP location: %v", err)

	if ipAddr == nil {
		return c.JSON(http.StatusNotFound, errors.New("IP address not found"))
	}

	return c.JSON(http.StatusOK, ipAddr)
}
