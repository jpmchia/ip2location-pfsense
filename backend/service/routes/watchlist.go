package routes

import (
	"errors"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/pfsense"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

var WatchList_GetRoute = "/api/watchlist"
var WatchList_PostItemRoute = "/api/watch"
var WatchList_GetItemRoute = "/api/watch"
var WatchList_DeleteItemRoute = "/api/watch"

func init() {
	util.LogDebug("[watchlist] Initialising widget service ...")

	util.Log("[watchlist] Watchlist add item endpoint: %s", WatchList_PostItemRoute)
	util.Log("[watchlist] Watchlist get item endpoint: %s", WatchList_GetItemRoute)
	util.Log("[watchlist] Watchlist delete item endpoint: %s", WatchList_DeleteItemRoute)
	util.Log("[watchlist] Watchlist get list endpoint: %s", WatchList_GetRoute)
}

// Responds to a request to add an IP address to the WatchList, storing the IP2Map entry in Redis and
// adding the IP address to the in memory list
func PostItemHandler(c echo.Context) error {
	ip_param := c.QueryParam("ip")
	ip2MapEntry := new(pfsense.Ip2Map)
	util.Log("[watchlist] Received request to add IP: %s\n", ip_param)

	if err := c.Bind(ip2MapEntry); err != nil {
		util.HandleError(err, "[watchlist] Failed to read the request body: %v", err)
		return c.String(http.StatusBadRequest, "[Watchlist] Bad Request. Failed to read the request body")
	}

	var ip = ip_param
	if pfsense.ActiveWatchList == nil {
		pfsense.ActiveWatchList = pfsense.NewWatchList()
	}

	if pfsense.ActiveWatchList.Contains(ip) {
		err := pfsense.ActiveWatchList.AddLogEntry(ip, *ip2MapEntry)
		util.HandleError(err, "[watchlist] Failed add the IP address to the list: %v", err)
	} else {
		pfsense.ActiveWatchList.Add(ip, *ip2MapEntry)
	}

	return c.JSON(http.StatusOK, pfsense.ActiveWatchList.GetWatchListDisplayItems())
}

// Responds to a request to get an item from the WatchList
func GetItemHandler(c echo.Context) error {
	ip_param := c.QueryParam("ip")
	util.Log("[watchlist] Received request to retrieve an item: %s\n", ip_param)

	if pfsense.ActiveWatchList == nil {
		pfsense.ActiveWatchList = pfsense.NewWatchList()
	}

	item, ok := pfsense.ActiveWatchList.Get(ip_param)
	if !ok {
		util.HandleError(errors.New("[watchList] Requested item was not found"), "[watchlist] WatchList item not found")
		return c.JSON(http.StatusNotFound, "Requested item was not found")
	}

	return c.JSON(http.StatusOK, item)
}

// Responds to a request to get the WatchList
func GetHandler(c echo.Context) error {
	detail := c.QueryParam("detail")

	if detail == "true" {
		return c.JSON(http.StatusOK, pfsense.ActiveWatchList)
	}

	return c.JSON(http.StatusOK, pfsense.ActiveWatchList.GetWatchListDisplayItems())
}

// Responds to a request to delete an item from the WatchList
func DeleteHandler(c echo.Context) error {
	ip_param := c.QueryParam("ip")
	util.Log("[watchlist] Received request to delete an item from the WatchList: %s\n", ip_param)

	pfsense.ActiveWatchList.Remove(ip_param)

	return c.JSON(http.StatusOK, pfsense.ActiveWatchList.GetWatchListDisplayItems())
}
