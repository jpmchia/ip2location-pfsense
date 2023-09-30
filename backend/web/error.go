package web

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/util"

	"github.com/labstack/echo/v4"
)

//go:embed error/*
var errorFiles embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

// Render implements echo.Renderer.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

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
	e.GET("/error/error.html.tmpl", echo.WrapHandler(http.StripPrefix("/error/", assetHandler)))
	return e
}

func ServeErrorTemplate(e *echo.Echo) *echo.Echo {
	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl")),
	}
	util.LogDebug("[webserve] Error template: %v", t.templates.Name())
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

		switch code = details.Code; code {
		case http.StatusNotFound:
			message = "Page not found"
		case http.StatusInternalServerError:
			message = "Internal server error"
		case http.StatusUnauthorized:
			message = "Not authorised"
		case http.StatusForbidden:
			message = "Forbidden"
		default:
			message = details.Message.(string)
		}
		log.Printf("CustomHTTPErrorHandler: %d : %v : %v", details.Code, details.Message, he.Internal)
	}

	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl")),
	}

	c.Echo().Renderer = t
	log.Printf("CustomHTTPErrorHandler: Rendering template with code %d and message %s", code, message)
	err = c.Render(http.StatusOK, "error.html.tmpl", map[string]interface{}{
		"code":    code,
		"message": message,
	})
	util.HandleError(err, "Unable to render error page")
}
