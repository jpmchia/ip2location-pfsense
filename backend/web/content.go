package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

//go:embed content/*
var contentFiles embed.FS

const contentPath = "content"

func embeddedContentFs() http.FileSystem {
	fsys, err := fs.Sub(contentFiles, contentPath)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func removePrefix(s string, prefix string) string {
	return s[len(prefix):]
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func ServeEmeddedContent(e *echo.Echo) *echo.Echo {
	contentFsHandler := http.FileServer(embeddedContentFs())

	files, err := getAllFilenames(&contentFiles)
	util.HandleError(err, "[web] Error loading content files")

	for _, file := range files {
		urlPath := removePrefix(file, "content/")
		e.GET(urlPath, echo.WrapHandler(contentFsHandler))
		util.LogDebug("[web] Registering: %v => %v\n", urlPath, file)
	}

	return e
}

func ContentHandler(c echo.Context) error {

	area := c.Param("area")

	switch area {
	case "dashboard":
		return DashboardHandler(c)

	case "watchlist":
		return WatchListHandler(c)

	case "ip2l":
		return Ip2lHandler(c)

	case "netstat":
		return NetstatHandler(c)

	case "service":
		return ServiceHandler(c)

	case "apikey":
		return ApiKeyHandler(c)

	case "settings":
		return SettingsHandler(c)

	case "about":
		return AboutHandler(c)

	case "help":
		return HelpHandler(c)

	case "login":
		return LoginHandler(c)

	case "logout":
		return LogoutHandler(c)

	default:
		return new(echo.HTTPError)

	}
}

type dashboard struct {
	Title string
}

func DashboardHandler(c echo.Context) error {

	var data dashboard
	data.Title = "IP2Location.io Backend service for pfSense"

	t := &TemplateRenderer{
		templates: template.Must(template.ParseFS(errorFiles, "error/error.html.templ")),
	}

	c.Echo().Renderer = t
	log.Printf("[web] DashboardHandler: Rendering template with title %s", data.Title)
	err := c.Render(http.StatusOK, "error.html.tmpl", map[string]interface{}{
		"title": data.Title,
	})
	util.HandleError(err, "[web] Unable to render error page")

	return c.String(http.StatusOK, "Dashboard")

}

func WatchListHandler(c echo.Context) error {
	return c.String(http.StatusOK, "WatchList")
}

func Ip2lHandler(c echo.Context) error {

	return c.String(http.StatusOK, "Ip2l")
}

func NetstatHandler(c echo.Context) error {

	return c.String(http.StatusOK, "Netstat")
}

func ServiceHandler(c echo.Context) error {

	return c.String(http.StatusOK, "Service")
}

func ApiKeyHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ApiKeys")
}

func SettingsHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Settings")
}

func AboutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "About")
}

func HelpHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Help")
}

func LoginHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}

func LogoutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Logout")
}
