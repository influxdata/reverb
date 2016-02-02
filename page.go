package reverb

import (
	"fmt"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/markbates/going/validate"
)

// Page is a value that can be used to passed data through
// the rendering of "github.com/flosch/pongo2" templates.
type Page struct {
	pongo2.Context
	Echo   *echo.Context
	layout string
}

// Layout returns the path to the current layout file.
// Example: "layouts/%s.html"
func (p Page) Layout() string {
	return fmt.Sprintf("layouts/%s.html", p.layout)
}

// SetTemplate set a value that defines the "yield"/"content"
// template.
func (p Page) SetTemplate(name string) {
	p.Context["_yield_template"] = name + ".html"
}

// Set adds the specificied value to the context. If a value
// of the same name already exists it will be overridden with the
// new value.
func (p Page) Set(key string, value interface{}) {
	p.Context[key] = value
}

// SetValidationErrors will add `validate.Errors` to the context
// using the key "validation_errors".
func (p Page) SetValidationErrors(verrs *validate.Errors) {
	p.Set("validation_errors", verrs)
}

// NewPage returns a new value of Page with defaults set. These
// defaults include the `layout` being set to "application" as well
// as the current request being set onto the `Context`.
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
