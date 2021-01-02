package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"someblocks/internal/app"
	"someblocks/internal/config"
	"someblocks/internal/controllers"
	"someblocks/internal/models"
	"someblocks/pkg/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	db := models.SetupAndMigrate(cfg)
	app := app.New(db)

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.SecretKey),
		csrf.Secure(cfg.Debug),
	)

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		csrfMiddleware,
		app.Session.LoadAndSave,
	)

	userService := models.UserService{DB: db}

	authController := controllers.NewAuthController(app, &userService)
	pageController := controllers.NewPageController(app)

	router.Get("/", pageController.PageIndex)
	router.Get("/page/{pageID}", app.LoginRequired(pageController.PageView))

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)
	router.Get("/auth/login", authController.Login)
	router.Post("/auth/login", authController.LoginPost)

	router.Post("/auth/logout", app.LoginRequired(authController.LogoutPost))

	router.Get("/auth/register", authController.Register)
	router.Post("/auth/register", authController.RegisterPost)

	// Setup static files /static route that will serve the static files from
	// from the ./static/ folder.
	filesDir := filepath.Join(utils.GetExecDir(), "static")
	route := "/static/"
	FileServer(router, route, filesDir)

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
