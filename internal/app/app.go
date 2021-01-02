package app

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"someblocks/internal/models"

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

	return &App{
		DB:      db,
		Render:  render,
		Session: sessionManager,
	}
}

func (app *App) GetCurrentUser(r *http.Request) *models.User {
	userId := app.Session.GetInt(r.Context(), "userId")
	if userId != 0 {
		us := models.UserService{DB: app.DB}
		user := us.GetById(userId)
		return user
	}
	return nil
}

func (app *App) LoginRequired(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.GetCurrentUser(r)
		if user == nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = WithCurrentUser(ctx, user)
		r = r.WithContext(ctx)
		f(w, r)
	}
}
