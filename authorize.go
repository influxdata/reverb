package reverb

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/markbates/going/imt"
)

// AuthorizeKey is the name of the key that will be used to
// set/get the "authorization" from the `echo.Context`.
var AuthorizeKey = "current_user"

// AuthorizeRedirectURL is the URL that the the `Authorize`
// middleware will redirect to for HTTP requests that require
// authorization, but aren't authorized.
var AuthorizeRedirectURL = "/"

// Authorize is an `echo` middleware that will check to make
// sure the `AuthorizeKey` variable is set in the `echo.Context`.
// If the key is not set it will redirect to the `AuthorizeRedirectURL`
// if the request is HTTP. If the request is of type `application/json`
// the a new `http.StatusForbidden` will be thrown.
func Authorize(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Get(AuthorizeKey) != nil {
			return handler(ctx)
		}
		ct := ctx.Request().Header().Get("Content-Type")
		if ct == imt.Application.JSON {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized Access")
		}
		return ctx.Redirect(302, AuthorizeRedirectURL)
	}
}
