package routes

import (
	"net/http"

	"github.com/jpmchia/ip2location-pfsense/backend/service/controller"
	"github.com/mikestefanello/pagoda/pkg/context"

	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	controller.Controller
}

func (e *ErrorHandler) Get(err error, ctx echo.Context) {
	if ctx.Response().Committed || context.IsCanceledError(err) {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code >= 500 {
		ctx.Logger().Error(err)
	} else {
		ctx.Logger().Info(err)
	}

	page := controller.NewPage(ctx)
	page.Layout = "main"
	page.Title = http.StatusText(code)
	page.Name = "error"
	page.StatusCode = code
	page.HTMX.Request.Enabled = false

	if err = e.RenderPage(ctx, page); err != nil {
		ctx.Logger().Error(err)
	}
}
