package reverb

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/markbates/going/validate"
	"github.com/markbates/reverb/sess"

	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	Data    map[string]interface{}
	lg      *log.Logger
	Session *sess.Session
	err     error
}

func (c *Context) Get(s string) interface{} {
	return c.Data[s]
}

func (c *Context) Set(s string, i interface{}) {
	c.Data[s] = i
}

func (c *Context) Logger() *log.Logger {
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

// SetTemplate set a value that defines the "yield"/"content"
// template.
func (c *Context) SetTemplate(name string) {
	c.Data["_yield_template"] = name + c.Extension()
}

// SetValidationErrors will add `validate.Errors` to the context
// using the key "validation_errors".
func (c *Context) SetValidationErrors(verrs *validate.Errors) {
	c.Set("validation_errors", verrs)
}

func NewContext(e echo.Context) *Context {
	c := &Context{
		Context: e,
		Data:    map[string]interface{}{},
		Session: sess.GetFromCtx(e),
		lg:      log.New(""),
	}

	c.lg.EnableColor()
	c.lg.SetFormat("${level}\t${time_rfc3339} ${long_file}:${line} : ${message}\n")

	c.SetExtension(".html")
	c.SetLayout("application")
	return c
}