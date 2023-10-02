package web

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jpmchia/ip2location-pfsense/backend/config"
	"github.com/jpmchia/ip2location-pfsense/backend/service/apikey"
	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

//go:embed content/* content/css/* content/fonts/* content/img/*
var contentFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

var templates map[string]*template.Template

const contentPath = "content"
const templatePath = "templates"
const errorPath = "error"

func init() {
	config.Configure()
	log.Printf("[web] Initialising templates")
	templates = make(map[string]*template.Template)
}

type MultiTemplateRenderer struct {
	templates map[string]*template.Template
}

func NoEscape(str string) template.HTML {
	return template.HTML(str)
}

// Render implements echo.Renderer.
func (t *MultiTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templ, ok := t.templates[name]
	if templ == nil || !ok {
		err := errors.New("Template not found: " + name)
		return err
	}

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return templ.ExecuteTemplate(w, "pfsense.html.tmpl", data)
}

func embeddedContentFs() http.FileSystem {
	fsys, err := fs.Sub(contentFiles, contentPath)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func embeddedTemplatesFs() http.FileSystem {
	fsys, err := fs.Sub(templateFiles, templatePath)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func embeddedErrorFs() http.FileSystem {
	fsys, err := fs.Sub(errorFiles, errorPath)
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func removePrefix(s string, prefix string) string {
	return s[len(prefix):]
}

func startsWithPrefix(path, prefix string) bool {
	// Clean the paths to remove any "..", ".", or multiple slashes etc.
	normalizedPath := filepath.Clean(path)
	normalizedPrefix := filepath.Clean(prefix)

	return strings.HasPrefix(normalizedPath, normalizedPrefix)
}

func getAllFilenames(efs ...*embed.FS) (files []string, err error) {

	for _, efs := range efs {

		if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			files = append(files, path)
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func ServeEmeddedContent(e *echo.Echo) *echo.Echo {
	contentFsHandler := http.FileServer(embeddedContentFs())
	templateFsHandler := http.FileServer(embeddedTemplatesFs())
	errorFsHandler := http.FileServer(embeddedErrorFs())

	files, err := getAllFilenames(&contentFiles, &templateFiles, &errorFiles)
	util.HandleError(err, "[web] Error loading content files")

	for _, file := range files {
		if startsWithPrefix(file, "error") {
			urlPath := removePrefix(file, "error/")
			e.GET(urlPath, echo.WrapHandler(errorFsHandler))
			util.LogDebug("[web] Registering: %v => %v\n", urlPath, file)
			continue
		}
		if startsWithPrefix(file, "templates") {
			urlPath := removePrefix(file, "templates/")
			e.GET(urlPath, echo.WrapHandler(templateFsHandler))
			util.LogDebug("[web] Registering: %v => %v\n", urlPath, file)
			continue
		}
		if startsWithPrefix(file, "content") {
			urlPath := removePrefix(file, "content/")
			e.GET(urlPath, echo.WrapHandler(contentFsHandler))
			util.LogDebug("[web] Registering: %v => %v\n", urlPath, file)
			continue
		}
	}
	return e
}

func ServeEmbeddedTemplates(e *echo.Echo) *echo.Echo {
	// First, create a new template with the required functions
	tmplWithFuncs := template.New("dashboard.html.tmpl").Funcs(template.FuncMap{
		"NoEscape": NoEscape,
	})

	ip2lTmplWithFuncs := template.New("ip2l.html.tmpl").Funcs(template.FuncMap{
		"NoEscape": NoEscape,
	})

	// Next, parse the templates with the defined functions
	templates["error.html.tmpl"] = template.Must(template.New("error.html.tmpl").ParseFS(errorFiles, "error/error.html.tmpl"))
	//templates["error.html.tmpl"] = template.Must(template.ParseFS(errorFiles, "error/error.html.tmpl"))
	templates["ip2l.html.tmpl"] = template.Must(ip2lTmplWithFuncs.ParseFS(templateFiles, "templates/ip2l.html.tmpl", "templates/layouts/pfsense.html.tmpl"))
	//templates["ip2l.html.tmpl"] = template.Must(template.ParseFS(templateFiles, "templates/ip2l.html.tmpl"))
	templates["watchlist.html.tmpl"] = template.Must(tmplWithFuncs.ParseFS(templateFiles, "templates/watchlist.html.tmpl", "templates/layouts/pfsense.html.tmpl"))
	//templates["watchlist.html.tmpl"] = template.Must(template.ParseFS(templateFiles, "templates/watchlist.html.tmpl"))
	templates["dashboard.html.tmpl"] = template.Must(tmplWithFuncs.ParseFS(templateFiles, "templates/dashboard.html.tmpl"))
	//templates["dashboard.html.tmpl"] = template.Must(template.ParseFS(templateFiles, "templates/dashboard.html.tmpl"))

	// templates["netstat.html.tmpl"] = template.Must(template.ParseFS(templateFiles, "templates/netstat.html.tmpl"))
	// templates["service.html.tmpl"] = template.Must(template.ParseFS(templateFiles, "templates/service.html.tmpl"))

	t := &MultiTemplateRenderer{
		templates: templates,
	}

	// t.templates["dashboard.html.tmpl"].Funcs(template.FuncMap{
	// 	"NoEscape": NoEscape,
	// })

	e.Renderer = t

	return e
}

func ContentHandler(c echo.Context) error {
	area := c.QueryParam("a")

	log.Printf("ContentHandler =====================> %v ===============> %v\n", area, c.Request().RequestURI)

	switch area {
	case "dashboard":
		util.Log("[web] DashboardHandler: %v", c.Request().RequestURI)
		return DashboardHandler(c)

	case "watchlist":
		util.Log("[web] WatchListHandler: %v", c.Request().RequestURI)
		return WatchListHandler(c)

	case "ip2l":
		util.Log("[web] Ip2lHandler: %v", c.Request().RequestURI)
		return Ip2lHandler(c)

	// case "netstat":
	// 	util.Log("[web] NetstatHandler: %v", c.Request().RequestURI)
	// 	return NetstatHandler(c)

	// case "service":
	// 	util.Log("[web] ServiceHandler: %v", c.Request().RequestURI)
	// 	return ServiceHandler(c)

	case "settings":
		util.Log("[web] SettingsHandler: %v", c.Request().RequestURI)
		return SettingsHandler(c)

	case "about":
		util.Log("[web] AboutHandler: %v", c.Request().RequestURI)
		return AboutHandler(c)

	case "help":
		util.Log("[web] HelpHandler: %v", c.Request().RequestURI)
		return HelpHandler(c)

	default:
		CustomHTTPErrorHandler(errors.New("page not found"), c)
		return nil
	}
}

func IncludeShaders(data map[string]interface{}) map[string]interface{} {

	data["shaderCode"] = `
	uniform sampler2D u_map_tex;

    varying float vOpacity;
    varying vec2 vUv;

    void main() {
        vec3 color = texture2D(u_map_tex, vUv).rgb;
        color -= .2 * length(gl_PointCoord.xy - vec2(.5));
        float dot = 1. - smoothstep(.38, .4, length(gl_PointCoord.xy - vec2(.5)));
        if (dot < 0.5) discard;
        gl_FragColor = vec4(color, dot * vOpacity);
    }`

	data["vertexCode"] = `
	uniform sampler2D u_map_tex;
	uniform float u_dot_size;
	uniform float u_time_since_click;
	uniform vec3 u_pointer;
	
	#define PI 3.14159265359
	
	varying float vOpacity;
	varying vec2 vUv;
	
	void main() {
	
		vUv = uv;
	
		// mask with world map
		float visibility = step(.2, texture2D(u_map_tex, uv).r);
		gl_PointSize = visibility * u_dot_size;
	
		// make back dots semi-transparent
		vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
		vOpacity = (1. / length(mvPosition.xyz) - .7);
		vOpacity = clamp(vOpacity, .03, 1.);
	
		// add ripple
		float t = u_time_since_click - .1;
		t = max(0., t);
		float max_amp = .15;
		float dist = 1. - .5 * length(position - u_pointer); // 0 .. 1
		float damping = 1. / (1. + 20. * t); // 1 .. 0
		float delta = max_amp * damping * sin(5. * t * (1. + 2. * dist) - PI);
		delta *= 1. - smoothstep(.8, 1., dist);
		vec3 pos = position;
		pos *= (1. + delta);
	
		gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.);
	}`

	return data
}

func DashboardHandler(c echo.Context) error {

	var data = make(map[string]interface{})
	var generatedKey = apikey.GenerateApiKey(c.RealIP(), 1800)

	data["Title"] = "IP2Location.io Backend service for pfSense"
	data["IPAddr"] = c.QueryParam("ip")
	data["RealIP"] = c.RealIP()
	data["Theme"] = c.QueryParam("theme")
	data["APIKey"] = generatedKey.Key
	data["APIKeyExpires"] = generatedKey.Expires
	data["Message"] = "Dashaboard!"
	data = IncludeShaders(data)

	log.Printf("[web] DashboardHandler: Rendering template with data %v", data)

	return c.Render(http.StatusOK, "dashboard.html.tmpl", data)
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
