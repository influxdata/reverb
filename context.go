package reverb

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/gommon/log"
	"github.com/markbates/going/validate"
	"github.com/markbates/reverb/sess"

	"github.com/labstack/echo"
	elog "github.com/labstack/echo/log"
)

type Context struct {
	echo.Context
	Data              map[string]interface{}
	lg                elog.Logger
	Session           *sess.Session
	RenderedTemplates []string
	err               error
	lock              *sync.Mutex
}

func (c *Context) RawRequest() *http.Request {
	return c.Request().(*standard.Request).Request
}

func (c *Context) RawResponse() http.ResponseWriter {
	return c.Response().(*standard.Response).ResponseWriter
}

func (c *Context) Get(s string) interface{} {
	return c.Data[s]
}

func (c *Context) Set(s string, i interface{}) {
	c.Data[s] = i
}

func (c *Context) Logger() elog.Logger {
	return c.lg
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Error(err error) {
	c.err = err
	c.Logger().Error(err)
}

func (c *Context) HandleError(err error) error {
	c.err = err
	c.Logger().Error(err)
	return err
}

func (c *Context) ContentType() string {
	return c.Request().Header().Get("Content-Type")
}

func (c *Context) Layout() string {
	return c.Data["_layout"].(string)
}

func (c *Context) SetLayout(s string) {
	c.Data["_layout"] = s
	if s != "" {
		c.Data["_layout"] = fmt.Sprintf("layouts/%s%s", s, c.Extension())
	}
}

func (c *Context) SetExtension(s string) {
	c.Set("extension", s)
}

func (c *Context) Extension() string {
	e, ok := c.Data["extension"]
	if ok {
		return e.(string)
	}
	return ".html"
}

// SetValidationErrors will add `validate.Errors` to the context
// using the key "validation_errors".
func (c *Context) SetValidationErrors(verrs *validate.Errors) {
	c.Set("validation_errors", verrs)
}

func (c *Context) LogRenderedTemplate(template string, fn func() error) error {
	now := time.Now()
	defer func() {
		stop := time.Now()
		end := stop.Sub(now)
		c.lock.Lock()
		c.RenderedTemplates = append(c.RenderedTemplates, template)
		c.lock.Unlock()
		lg := c.Get("lg").(*Logger)
		lg.AddDurations("Views", end)
		lg.Printf("  Rendered %s %s", template, end)
	}()
	return fn()
}

func NewContext(e echo.Context) *Context {
	lg := log.New("")
	c := &Context{
		Context:           e,
		Data:              map[string]interface{}{},
		Session:           sess.GetFromCtx(e),
		RenderedTemplates: []string{},
		lock:              &sync.Mutex{},
		lg:                lg,
	}

	lg.EnableColor()
	// lg.SetFormat("${level}\t${time_rfc3339} ${long_file}:${line} : ${message}\n")

	c.SetExtension(".html")
	c.SetLayout("application")
	return c
}
