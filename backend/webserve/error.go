package webserve

import (
	"embed"
	"html/template"
	"io/fs"
	"ip2location-pfsense/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed error/*
var errorFiles embed.FS

// Gets the embedded file system
func getErrorFileSystem() http.FileSystem {
	fsys, err := fs.Sub(errorFiles, "error")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

// ServeEmbedded serves the embedded files
func ServeEmbeddedErrorFiles(e *echo.Echo) *echo.Echo {
	assetHandler := http.FileServer(getErrorFileSystem())
	// e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/error/error.html.tmpl", echo.WrapHandler(http.StripPrefix("/error/", assetHandler)))
	return e
}

func ServeErrorTemplate(e *echo.Echo) *echo.Echo {
	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl")),
	}

	util.LogDebug("Template: %v", t)

	e.Renderer = t

	return e
}

// CustomHTTPErrorHandler is a custom error handler that renders the error page
// from the embedded file system or renders using a template, also embedded,
// depending on the error code
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var message string
	code := http.StatusInternalServerError
	c.Logger().Error(err)

	if he, ok := err.(*echo.HTTPError); ok {
		details := err.(*echo.HTTPError)
		code = details.Code
		message = details.Message.(string)
		if code == http.StatusNotFound {
			message = "Page not found"
		}
		if code == http.StatusInternalServerError {
			message = "Internal server error"
		}
		if code == http.StatusUnauthorized {
			message = "Unauthorized"
		}
		if code == http.StatusForbidden {
			message = "Forbidden"
		}
		util.LogDebug("CustomHTTPErrorHandler: %d : %v : %v", details.Code, details.Message, he.Internal)
	}

	// LogDebug("CustomHTTPErrorHandler: %d:%v", details.Code, details.Message)
	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl")),
	}

	c.Echo().Renderer = t
	err = c.Render(http.StatusOK, "error.html.tmpl", map[string]interface{}{
		"code":    code,
		"message": message,
	})
	util.HandleError(err, "Unable to render error page")
}
