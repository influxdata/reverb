package reverb

import (
	"fmt"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/markbates/going/validate"
)

type Page struct {
	pongo2.Context
	Echo   *echo.Context
	layout string
}

func (p Page) Layout() string {
	return fmt.Sprintf("layouts/%s.html", p.layout)
}

func (p Page) SetTemplate(name string) {
	p.Context["_yield_template"] = name + ".html"
}

func (p Page) Set(key string, value interface{}) {
	p.Context[key] = value
}

func (p Page) SetValidationErrors(verrs *validate.Errors) {
	p.Context["validation_errors"] = verrs
}

func NewPage(ctx *echo.Context) Page {
	p := Page{
		Echo: ctx,
		Context: pongo2.Context{
			"request": ctx.Request(),
		},
		layout: "application",
	}
	return p
}
