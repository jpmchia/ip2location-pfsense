package webserve

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"ip2location-pfsense/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed static/index.html.tmpl
var embededFiles embed.FS

func getFileSystem() http.FileSystem {
	util.LogDebug("[webserve] Serving up embedded files ...")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func ServeEmbedded(e *echo.Echo) *echo.Echo {
	assetHandler := http.FileServer(getFileSystem())
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	return e
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func ServeRenderTemplate(e *echo.Echo) *echo.Echo {
	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(embededFiles, "static/index.html.tmpl")),
	}
	util.LogDebug("Template: %v", t)

	e.Renderer = t

	e.GET("/index.html", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html.tmpl", map[string]interface{}{
			"code":    "ip2location",
			"message": "Hello",
		})
	})

	return e
}
