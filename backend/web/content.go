package web

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

	files, err := getAllFilenames(&contentFiles)
	util.HandleError(err, "[web] Error loading content files")

	for _, file := range files {
		urlPath := removePrefix(file, "content/")
		e.GET(urlPath, echo.WrapHandler(contentFsHandler))
		util.LogDebug("[web] Registering: %v => %v\n", urlPath, file)
	}

	return e
}
