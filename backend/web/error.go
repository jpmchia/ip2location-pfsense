package web

// type TemplateRenderer struct {
// 	templates *template.Template
// }

// Render implements echo.Renderer.
// func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

// 	if viewContext, isMap := data.(map[string]interface{}); isMap {
// 		viewContext["reverse"] = c.Echo().Reverse
// 	}

// 	return t.templates.ExecuteTemplate(w, name, data)
// }

// Gets the embedded file system
// func getErrorFileSystem() http.FileSystem {
// 	fsys, err := fs.Sub(errorFiles, "error")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return http.FS(fsys)
// }

// ServeEmbedded serves the embedded files
// func ServeEmbeddedErrorFiles(e *echo.Echo) *echo.Echo {
// 	assetHandler := http.FileServer(getErrorFileSystem())
// 	e.GET("/error/error.html.tmpl", echo.WrapHandler(http.StripPrefix("/error/", assetHandler)))
// 	return e
// }

// func ServeErrorTemplate(e *echo.Echo) *echo.Echo {
// 	t := &TemplateRenderer{
// 		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl")),
// 	}
// 	util.LogDebug("[webserve] Error template: %v", t.templates.Name())
// 	e.Renderer = t
// 	return e
// }
