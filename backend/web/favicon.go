package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed favicon/*
var faviconFiles embed.FS

const favIconPath = "favicon"
const siteManifest = "site.webmanifest"
const favIcon = "favicon.ico"
const favIcon16 = "favicon-16x16.png"
const favIcon32 = "favicon-32x32.png"

func embeddedFaviconHandler() http.FileSystem {
	fsys, err := fs.Sub(faviconFiles, favIconPath)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func ServeFavIcons(e *echo.Echo) *echo.Echo {
	faviconFsHandler := http.FileServer(embeddedFaviconHandler())
	e.GET(siteManifest, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon16, echo.WrapHandler(faviconFsHandler))
	e.GET(favIcon32, echo.WrapHandler(faviconFsHandler))
	return e
}
