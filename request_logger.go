package reverb

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/flosch/go-humanize"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			lg := NewLogger(c)
			c.Set("lg", lg)

			start := time.Now()

			req := c.Request()
			res := c.Response()

			lg.Printf("Started %s \"%s\" for %s %s", req.Method, path(req), remoteAddr(req), start)

			err := h(c)
			if err != nil {
				lg.Printf("  Error: %s", err)
				c.Error(err)
			}

			stop := time.Now()
			size := res.Size()

			lg.AddExtras(fmt.Sprintf("Size: %s", humanize.Bytes(uint64(size))))
			ds := lg.Durations.String()
			if ds != "" {
				lg.AddExtras(ds)
			}

			lg.Printf("Completed %s in %s%s", code(res), stop.Sub(start), lg.Extras)
			return nil
		}
	}
}

func remoteAddr(req *http.Request) string {
	var remoteAddr string
	if ip := req.Header.Get(echo.XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(echo.XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr = req.RemoteAddr
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	return remoteAddr
}

func path(req *http.Request) string {
	path := req.URL.Path

	if path == "" {
		path = "/"
	}
	return path
}

func code(res *echo.Response) string {
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
