package app

import (
	"html/template"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type App struct {
	render *render.Render
	DB *sqlx.DB
}

func New() *App {
	app := App{}

	// Setup "Template Engine" AKA renderer
	app.render = render.New(render.Options{
		RenderPartialsWithoutPrefix: true,
		IsDevelopment:               viper.GetBool("debug"),
		Directory:                   "templates",
		Layout:                      "base",
		Extensions:                  []string{".html"},
		Funcs: []template.FuncMap{
			// Will be overriden in AppContext.HTML to add a CSRF Field
			template.FuncMap{"csrfField": func() string {
				return ""
			},
			},
		},
	})

	return &app
}
