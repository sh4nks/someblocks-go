package core

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

// This holds the app context
type AppContext struct {
	Render *render.Render
	DB     *sqlx.DB
	Port   string
	Host   string
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

func (c *AppContext) JSON(w http.ResponseWriter, status int, v interface{}) {
	c.Render.JSON(w, status, v)
}

func (c *AppContext) Text(w http.ResponseWriter, status int, v string) {
	c.Render.Text(w, status, v)
}

func (c *AppContext) HTML(w http.ResponseWriter, r *http.Request, status int, tmpl string, data interface{}) {
	csrfField := csrf.TemplateField(r)
	htmlOpts := render.HTMLOptions{
		Funcs: template.FuncMap{
			"csrfField": func() template.HTML {
				return csrfField
			},
			"testFunc": func() string {
					return "My custom function"
				},
		},
	}
	c.Render.HTML(w, status, tmpl, data, htmlOpts)
}
