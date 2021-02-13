package app

import (
	"encoding/gob"
	"html/template"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

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
			// these functions will be overriden in app.HTML due to them
			// needing a http.Request object
			template.FuncMap{
				"csrfField": func() string {
					return ""
				},
				"getFlashedMessages": func() *Flash {
					return &Flash{}
				},
				"isActive": func() string {
					return ""
				},
			},
		},
	})

	sessionManager := &SessionManager{
		*scs.New(),
	}
	sessionManager.Store = sqlite3store.New(db.DB)

	return &App{
		DB:      db,
		Render:  render,
		Session: sessionManager,
	}
}
