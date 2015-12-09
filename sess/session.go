package sess

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

const sessionName = "_influxdata_enterprise_session"

type Session struct {
	Session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

func (s *Session) Save() error {
	return s.Session.Save(s.req, s.res)
}

func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

func Get(r *http.Request, w http.ResponseWriter) *Session {
	session, _ := store.Get(r, sessionName)
	return &Session{
		Session: session,
		req:     r,
		res:     w,
	}
}

func GetFromCtx(ctx *echo.Context) *Session {
	return Get(ctx.Request(), ctx.Response())
}
