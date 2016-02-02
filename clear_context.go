package reverb

import (
	"github.com/gorilla/context"
	"github.com/labstack/echo"
)

// ClearContext is an `echo` middle that clears the `echo.Context`
// at the end of a request. You should DEFINITELY add this to your
// middleware stack! If you don't, you run the risk of memory bloat!
func ClearContext(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		defer context.Clear(ctx.Request())
		return handler(ctx)
	}
}
