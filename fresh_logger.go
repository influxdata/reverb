package reverb

import (
	"net/http"

	"github.com/pilu/fresh/runner/runnerutils"
)

func FreshLogger(app http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if runnerutils.HasErrors() {
			runnerutils.RenderError(res)
		} else {
			app.ServeHTTP(res, req)
		}
	}
}
