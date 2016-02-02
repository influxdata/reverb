package reverb

import (
	"net/http"

	"github.com/pilu/fresh/runner/runnerutils"
)

// FreshLogger is an `echo` middleware for serving up
// error pages using the "github.com/pilu/fresh/runner/runnerutils"
// package.
func FreshLogger(app http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if runnerutils.HasErrors() {
			runnerutils.RenderError(res)
		} else {
			app.ServeHTTP(res, req)
		}
	}
}
