package reverb

import (
	"github.com/gorilla/context"
	"github.com/labstack/echo"
)

func ClearContext(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		defer context.Clear(ctx.Request())
		return handler(ctx)
	}
}
