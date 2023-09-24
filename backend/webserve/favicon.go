package webserve

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

//go:embed static/favicon/*
var faviconFiles embed.FS

const favIconPath = "/static/favicon"
const siteManifest = "site.webmanifest"
const favIcon = "favicon.ico"
const favIcon16 = "favicon-16x16.png"
const favIcon32 = "favicon-32x32.png"

func embeddedFaviconHandler() http.FileSystem {
	fsys, err := fs.Sub(faviconFiles, "static/favicon")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func ServeFavIcons(e *echo.Echo) *echo.Echo {
	err := LoadTemplates()
	util.HandleError(err, "LoadTemplates")
	faviconFsHandler := http.FileServer(embeddedFaviconHandler())
	e.GET(siteManifest, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon16, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon32, echo.WrapHandler(faviconFsHandler))
	return e
}
