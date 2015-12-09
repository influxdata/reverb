package reverb

import (
	"net/http"

	"github.com/markbates/going/defaults"
)

var HeaderHTTPMethodOverride = "X-HTTP-Method-Override"
var ParamHTTPMethodOverride = "_method"

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
