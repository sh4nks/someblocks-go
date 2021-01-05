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
			// Will be overriden in "(app *App) HTML()" to add a CSRF Field and
			// a display the flashed messages
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
	sessionManager.Store = sqlite3store.New(db.DB)

	return &App{
		DB:      db,
		Render:  render,
		Session: sessionManager,
	}
}
