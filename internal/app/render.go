package app

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/unrolled/render"
)

func (app *App) JSON(w http.ResponseWriter, status int, v interface{}) {
	app.Render.JSON(w, status, v)
}

func (app *App) Text(w http.ResponseWriter, status int, v string) {
	app.Render.Text(w, status, v)
}

func (app *App) HTML(w http.ResponseWriter, r *http.Request, status int, tmpl string, data interface{}) {
	htmlOpts := render.HTMLOptions{
		Funcs: template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"getFlashedMessages": func() *Flash {
				flash := app.Session.GetFlash(r.Context())
				if flash != nil {
					return flash
				}
				return nil
			},
		},
	}
	app.Render.HTML(w, status, tmpl, data, htmlOpts)
}
