package reverb

import "github.com/labstack/echo"

// DefautlContentType is an `echo` middleware that will
// the request "Content-Type" header to the specified type
// if the header is not already set.
func DefaultContentType(s string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header().Get("Content-Type") == "" {
				c.Request().Header().Set("Content-Type", s)
			}
			return h(c)
		}
	}
}
