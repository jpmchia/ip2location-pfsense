package webserve

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed static/*
var staticFiles embed.FS

// // var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})
// var staticFilesHandler = echo.WrapHandler(http.FileServer(http.FS(staticFiles)))

// // Gets the embedded file system
func getStaticFilesFs() http.FileSystem {
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func ServeStaticFiles(e *echo.Echo) *echo.Echo {
	staticFilesHandler := http.FileServer(getStaticFilesFs())
	e.GET("/", echo.WrapHandler(staticFilesHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticFilesHandler)))
	return e
}

// Middleware to rewrite the path to the static files
// func contentRewrite(next echo.HandlerFunc) echo.HandlerFunc {
// 	next = middleware.Rewrite(map[string]string{"/*": "/static/$1"})
// 	return next
// }

// func ServeStaticFiles(e *echo.Echo) *echo.Echo {
// 	err := LoadTemplates()
// 	util.HandleError(err, "LoadTemplates")
// 	e.Pre(middleware.RewriteWithConfig(map[string]string{"/*": "/static/$1"))
// 	e.GET("/*", echo.WrapHandler(staticFilesHandler), contentRewrite)
// 	return e
// }

// var staticFilesHandler = echo.WrapHandler(http.FileServer(http.FS(staticFiles)))
// var staticFilesRewrite = echo.RewriteWithConfig(map[string]string{"/*": "/public/$1"})
