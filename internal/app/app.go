package app

import (
	"encoding/gob"
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type App struct {
	DB      *sqlx.DB
	Session *SessionManager
	Render  *render.Render
}

func init() {
	gob.Register(Flash{})
}

func New(db *sqlx.DB) *App {
	// Setup "Template Engine" AKA renderer
	render := render.New(render.Options{
		RenderPartialsWithoutPrefix: true,
		IsDevelopment:               viper.GetBool("debug"),
		Directory:                   "templates",
		Layout:                      "base",
		Extensions:                  []string{".html"},
		Funcs: []template.FuncMap{
			// Will be overriden in AppContext.HTML to add a CSRF Field
			template.FuncMap{
				"csrfField": func() string {
					return ""
				},
				"getFlashedMessages": func() *Flash {
					return &Flash{}
				},
			},
		},
	})

	sessionManager := &SessionManager{
		*scs.New(),
	}

	return &App{
		DB:      db,
		Render:  render,
		Session: sessionManager,
	}
}
