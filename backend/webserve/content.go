package webserve

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/util"
	"github.com/labstack/echo/v4"
)

//go:embed content/*
var contentFiles embed.FS

const contentPath = "content"

func embeddedContentHandler() http.FileSystem {
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
	contentFsHandler := http.FileServer(embeddedContentHandler())
	//e.GET("ip2location.html", echo.WrapHandler(contentFsHandler))
	// e.GET("bundle.js", echo.WrapHandler(contentFsHandler))
	// e.GET("bundle.js.map", echo.WrapHandler(contentFsHandler))
	// e.GET("style.css", echo.WrapHandler(contentFsHandler))
	// e.GET("css/style.css", echo.WrapHandler(contentFsHandler))

	files, err := getAllFilenames(&contentFiles)
	util.HandleError(err, "[webserve] Error loading content files")

	for _, file := range files {
		urlPath := removePrefix(file, "content/")
		e.GET(urlPath, echo.WrapHandler(contentFsHandler))
		util.LogDebug("[webserve] Registering: %v => %v\n", urlPath, file)
	}

	return e
}

// var files map[string]string

// const contentPath = "content/*"

// func getAllFilenames(efs *embed.FS) (files map[string]string, err error) {
// 	files = make(map[string]string)
// 	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
// 		if d.IsDir() {
// 			return nil
// 		}
// 		file = d.Name()
// 		//files = append(files, path)
// 		urlPath := removePrefix(files, "content/")
// 		files[urlPath] = files
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}
// 	return files, nil
// }

// func contentFileHandler() http.FileSystem {
// 	fsys, err := fs.Sub(contentFiles, "content")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return http.FS(fsys)
// }

// func ServeEmbeddedContent(e *echo.Echo) *echo.Echo {
// 	embeddedFileHandler := http.FileServer(contentFileHandler())

// 	files, err := getAllFilenames(&contentFiles)

// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, file := range files {

// 		e.GET(file, echo.WrapHandler(embeddedFileHandler))
// 		fmt.Printf("Registering: %v => %v\n", file, file)
// 	}
// 	e.GET("/", echo.WrapHandler(embeddedFileHandler))
// 	e.GET("/content/*", echo.WrapHandler(http.StripPrefix("/content/", embeddedFileHandler)))
// 	return e
// }
