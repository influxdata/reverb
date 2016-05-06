package reverb_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/markbates/going/willy"
	"github.com/markbates/reverb"
	"github.com/stretchr/testify/require"
)

func Test_DefaultContentType(t *testing.T) {
	r := require.New(t)

	e := echo.New()
	e.Use(reverb.DefaultContentType("foo/bar"))
	e.Get("/", func(c echo.Context) error {
		return c.String(200, c.Request().Header().Get("Content-Type"))
	})

	w := willy.New(reverb.StandardHandler(e))
	res := w.Request("/").Get()
	r.Equal("foo/bar", res.Body.String())
}
