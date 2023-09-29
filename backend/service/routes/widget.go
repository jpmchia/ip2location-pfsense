package routes

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

var Widget_PostRoute string
var Widget_GetRoute string

func init() {
	util.LogDebug("[widget] Initialising widget service ...")

	Widget_PostRoute = "/widget/configuration"
	Widget_GetRoute = "/widget/controls"

	util.Log("[widget] Widget configuration service endpoint: %s", Widget_PostRoute)
	util.Log("[widget] Widget content service endpoint: %s", Widget_GetRoute)
}

func PostWidgetHandler(c echo.Context) error {
	return c.String(http.StatusNoContent, "Not Implemented")
}

func GetWidgetHandler(c echo.Context) error {
	return c.String(http.StatusNoContent, "Not Implemented")
}
