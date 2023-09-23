package webserve

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

const (
	layoutsDir  = "templates/layouts"
	templateDir = "templates"
	extension   = "/*.html.tmpl"
)

var (
	//go:embed templates/* templates/layouts/*
	files     embed.FS
	templates map[string]*template.Template
)

//go:embed static/*
var staticFiles embed.FS

const testPage = "template.html.tmpl"

type TemplateRenderer struct {
	templates *template.Template
}

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templateDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if !tmpl.IsDir() {
			continue
		}
		tmplName := tmpl.Name()
		util.LogDebug("Loading template: %s", tmplName)
		templates[tmplName], err = template.ParseFS(files, templateDir+"/"+tmplName+extension)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestPage(w http.ResponseWriter, r *http.Request) {
	t, ok := templates[testPage]
	if !ok {
		log.Printf("template %s not found", testPage)
		return
	}

	data := make(map[string]interface{})
	data["code"] = "This is a test page"
	data["message"] = "Welcome to the test page"

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func ServeEmbedded(e *echo.Echo) *echo.Echo {
	err := LoadTemplates()
	util.HandleError(err, "LoadTemplates")
	//assetHandler := http.FileServer(getFileSystem())
	// e.GET("/", echo.WrapHandler(assetHandler))
	//e.GET("/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	//e.GET("/", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	return e
}

// func getFileSystem() http.FileSystem {
// 	util.LogDebug("[webserve] Serving up embedded static files ...")
// 	fsys, err := fs.Sub(staticFiles, "static")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return http.FS(fsys)
// }

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func ServeRenderTemplate(e *echo.Echo) *echo.Echo {
	// t := &TemplateRenderer{
	// 	templates: template.Must(template.ParseFS(staticFiles, "static/template.html.tmpl")),
	// }
	r := http.NewServeMux()
	r.HandleFunc("/test.html", TestPage)

	// util.LogDebug("Template: %v", t)
	// e.Renderer = t
	// e.GET("/test.html", func(c echo.Context) error {
	// 	return c.Render(http.StatusOK, "template.html.tmpl", map[string]interface{}{
	// 		"content": "<a href=/login.html>login</a>",
	// 	})
	// })
	return e
}
