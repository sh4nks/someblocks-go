package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"someblocks/app"
	"someblocks/config"
	"someblocks/controllers"
	"someblocks/database"
	"someblocks/middleware"
	"someblocks/models"
	"someblocks/utils"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

// https://github.com/infomark-org/infomark/blob/62367a65aadf3e38f7ee9cfb1401180d04374b52/api/server.go
// https://github.com/infomark-org/infomark/blob/62367a65aadf3e38f7ee9cfb1401180d04374b52/api/app/api.go#L244
// https://gitlab.com/joncalhoun/lenslocked.com/-/blob/master/main.go

type Server struct {
	App    *app.App
	Router *chi.Mux
	Config *config.Config
}

func New(cfg *config.Config) *Server {
	db, err := database.SetupAndMigrate(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't setup database")
		return nil
	}

	app := app.New(db)

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.SecretKey),
		csrf.Secure(cfg.Debug),
		csrf.Path("/"),
	)

	router := chi.NewRouter()
	router.Use(
		chiMiddleware.RequestID,
		chiMiddleware.Logger,
		chiMiddleware.RedirectSlashes,
		chiMiddleware.Recoverer,
		csrfMiddleware,
		app.Session.LoadAndSave,
	)

	userService := models.UserService{DB: db}

	authController := controllers.NewAuthController(app, &userService)
	userController := controllers.NewUserController(app, &userService)
	pageController := controllers.NewPageController(app)

	userMw := middleware.NewUserMiddleware(app, &userService)
	router.Use(userMw.CurrentUser)

	router.Get("/", pageController.PageIndex)
	router.With(userMw.LoginRequired).Get("/page/{pageID}", pageController.PageView)

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)
	router.Get("/auth/login", authController.Login)
	router.Post("/auth/login", authController.LoginPost)

	router.With(userMw.LoginRequired).Post("/auth/logout", authController.LogoutPost)

	router.Get("/auth/register", authController.Register)
	router.Post("/auth/register", authController.RegisterPost)

	router.Mount("/user", userController.Routes(userMw))

	// Setup static files /static route that will serve the static files from
	// from the ./static/ folder.
	filesDir := filepath.Join(utils.GetExecDir(), "ui", "dist")
	FileServer(router, "/static/", filesDir)

	return &Server{
		Config: cfg,
		Router: router,
		App:    app,
	}
}

func (srv *Server) Start() {
	addr := fmt.Sprintf("%s:%d", srv.Config.Web.Host, srv.Config.Web.Port)
	log.Info().Msgf("Running on http://%s/ (Press CTRL+C to quit)", addr)

	http.ListenAndServe(addr, srv.Router)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r *chi.Mux, endpoint string, filesDir string) {
	if strings.ContainsAny(endpoint, "{}*") {
		panic("static file server does not permit any URL parameters.")
	}

	// Check if the path ends with '/' - if not return 404
	if endpoint != "/" && endpoint[len(endpoint)-1] != '/' {
		r.Get(endpoint, r.NotFoundHandler())
	}

	fs := http.StripPrefix(endpoint, http.FileServer(http.Dir(filesDir)))
	r.Get(endpoint+"*", func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, endpoint, "/", 1)
		if _, err := os.Stat(filesDir + file); os.IsNotExist(err) {
			http.ServeFile(w, r, file)
			return
		}
		fs.ServeHTTP(w, r)
	})
}
