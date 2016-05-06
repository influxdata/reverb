package reverb

import (
	"net/http"

	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

func StandardHandler(h engine.Handler) http.Handler {
	s := standard.New(":9999")
	s.SetHandler(h)
	return s.Handler
}
