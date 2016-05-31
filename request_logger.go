package reverb

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/flosch/go-humanize"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/gommon/color"
)

var AssetsPath = regexp.MustCompile("^/assets/.+")

// RequestLogger is an `echo` middleware that wraps a request
// and nicely formats it using the `reverb.Logger`.
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request().(*standard.Request).Request
		res := c.Response()
		path := path(c)

		lg := NewLogger(c)
		c.Set("lg", lg)

		// don't log assets
		if AssetsPath.MatchString(path) {
			lg.SetOutput(&bytes.Buffer{})
		}

		start := time.Now()

		lg.Printf("Started %s \"%s\" for %s %s", req.Method, path, remoteAddr(req), start)

		err := next(c)

		stop := time.Now()
		size := res.Size()

		lg.AddExtras(fmt.Sprintf("Size: %s", humanize.Bytes(uint64(size))))

		ds := lg.Durations.String()
		if ds != "" {
			lg.AddExtras(ds)
		}

		lg.Printf("Completed %s in %s%s", code(res), stop.Sub(start), lg.Extras)
		return err
	}

}

func remoteAddr(req *http.Request) string {
	var remoteAddr string
	if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr = req.RemoteAddr
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	return remoteAddr
}

func path(ctx echo.Context) string {
	path := ctx.Request().URL().Path()

	if path == "" {
		path = "/"
	}
	return path
}

func code(res engine.Response) string {
	n := res.Status()
	code := color.Green(n)
	switch {
	case n >= 500:
		code = color.Red(n)
	case n >= 400:
		code = color.Yellow(n)
	case n >= 300:
		code = color.Cyan(n)
	}
	return code
}
