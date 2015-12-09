package reverb

import "github.com/labstack/echo"

func ErrorCatcher(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := h(c)
		if err != nil {
			c.Error(err)
		}
		return nil
	}
}
