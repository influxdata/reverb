package reverb

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/markbates/going/imt"
)

var AuthorizeKey = "current_user"
var AuthorizeRedirectURL = "/"

func Authorize(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		if ctx.Get(AuthorizeKey) != nil {
			return handler(ctx)
		}
		ct := ctx.Request().Header.Get("Content-Type")
		if ct == imt.Application.JSON {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized Access")
		}
		return ctx.Redirect(302, AuthorizeRedirectURL)
	}
}
