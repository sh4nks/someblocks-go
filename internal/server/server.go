package server

import (
	"fmt"
	"html/template"
	"io"

	"someblocks/internal/app"
	"someblocks/internal/config"
	"someblocks/internal/controllers"
	"someblocks/internal/models"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"
)

// https://github.com/infomark-org/infomark/blob/62367a65aadf3e38f7ee9cfb1401180d04374b52/api/server.go
// https://github.com/infomark-org/infomark/blob/62367a65aadf3e38f7ee9cfb1401180d04374b52/api/app/api.go#L244
// https://gitlab.com/joncalhoun/lenslocked.com/-/blob/master/main.go


type Server struct {
	App    *app.App
	Config *config.Config
	Router *echo.Echo
}

type RenderWrapper struct { // We need to wrap the renderer because we need a different signature for echo.
	rnd *render.Render
}

func (r *RenderWrapper) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	htmlOpts := render.HTMLOptions{
		Funcs: template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(c.Request())
			},
		},
	}

	return r.rnd.HTML(w, 0, name, data, htmlOpts) // The zero status code is overwritten by echo.
}


func New(cfg *config.Config) *Server {
	db := models.SetupAndMigrate(cfg)

	appData := app.New()
	appData.DB = db

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.SecretKey),
		csrf.Secure(cfg.Debug),
	)

	renderer := &RenderWrapper{
		render.New(
			render.Options{
				RenderPartialsWithoutPrefix: true,
				IsDevelopment:               cfg.Debug,
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
			})}

	// Echo instance
	e := echo.New()

	// Template engine
	e.Renderer = renderer

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echo.WrapMiddleware(csrfMiddleware))

	pageController := controllers.NewPageController(appData)
	authController := controllers.NewAuthController(appData)

	e.GET("/", pageController.PageIndex)
	e.GET("/page/:pageID", pageController.PageView)

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)

	e.GET("/auth/login", authController.Login)
	e.POST("/auth/logout", authController.Logout)

	e.Static("/static", "static")

	return &Server{
		Config: cfg,
		Router: e,
		App:    appData,
	}
}

func (srv *Server) Start() {
	addr := fmt.Sprintf("%s:%d", srv.Config.Web.Host, srv.Config.Web.Port)
	log.Info().Msgf("Running on http://%s/ (Press CTRL+C to quit)", addr)
	srv.Router.Start(addr)
}
