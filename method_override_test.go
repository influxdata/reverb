package reverb_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/markbates/going/willy"
	"github.com/markbates/reverb"
	"github.com/stretchr/testify/require"
)

func Test_MethodOverride(t *testing.T) {
	r := require.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, req.Method)
	})

	w := willy.New(reverb.MethodOverride(mux))
	res := w.Request("/").Get()
	r.Equal("GET", res.Body.String())

	res = w.Request("/").Post(url.Values{})
	r.Equal("POST", res.Body.String())

	for _, v := range []string{"PUT", "PATCH", "DELETE"} {
		res := w.Request("/").Post(url.Values{reverb.ParamHTTPMethodOverride: []string{v}})
		r.Equal(v, res.Body.String())
	}
}
