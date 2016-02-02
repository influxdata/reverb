package reverb

import (
	"net/http"

	"github.com/markbates/going/defaults"
)

// HeaderHTTPMethodOverride is the key of the HTTP header used
// to override the request method in `MethodOverride`.
var HeaderHTTPMethodOverride = "X-HTTP-Method-Override"

// ParamHTTPMethodOverride is the form value used to override the
// request method in `MethodOverride`.
var ParamHTTPMethodOverride = "_method"

// MethodOverride allows requests that don't support HTTP methods such
// as "PUT", "PATCH", "DELETE", to use either a header, `HeaderHTTPMethodOverride`,
// or a form value, `ParamHTTPMethodOverride`, to change the request
// method.
func MethodOverride(app http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			_method := defaults.String(req.FormValue(ParamHTTPMethodOverride), req.Header.Get(HeaderHTTPMethodOverride))
			switch _method {
			case "PUT", "PATCH", "DELETE":
				req.Method = _method
			}
		}
		app.ServeHTTP(res, req)
	}
}
