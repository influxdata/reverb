package reverb_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/markbates/going/willy"
	"github.com/markbates/reverb"
	"github.com/stretchr/testify/require"
)

func Test_Authorize(t *testing.T) {
	r := require.New(t)

	h := func(ctx echo.Context) error {
		return ctx.String(200, "HI!")
	}

	authorizer := func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(reverb.AuthorizeKey, "some user")
			return handler(ctx)
		}
	}

	e := echo.New()
	e.Get("/unauthorized", reverb.Authorize(h))
	e.Get("/authorized", authorizer(reverb.Authorize(h)))

	w := willy.New(reverb.StandardHandler(e))
	res := w.Request("/unauthorized").Get()
	r.Equal(302, res.Code)
	r.Equal(reverb.AuthorizeRedirectURL, res.Location())

	res = w.Request("/authorized").Get()
	r.Equal(200, res.Code)
	r.Equal("HI!", res.Body.String())
}
