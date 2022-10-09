package app

import (
	"encoding/gob"
	"html/template"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
	"gorm.io/gorm"
)

type App struct {
	DB      *gorm.DB
	Session *SessionManager
	Render  *render.Render
}

func init() {
	gob.Register(Flash{})
}

func New(db *gorm.DB) *App {
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

	store, err := gormstore.New(db)
	if err != nil {
		log.Fatal().Msg("Error while migrating the session store")
	}
	sessionManager.Store = store

	return &App{
		DB:      db,
		Render:  render,
		Session: sessionManager,
	}
}
