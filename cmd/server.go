package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"someblocks/core"
	"someblocks/routes/page"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateApp() http.Handler {
	dir, _ := os.Getwd()
	dir = filepath.Clean(filepath.Join(dir, ".."))
	db, err := sqlx.Open("sqlite3", filepath.Join(dir, "sqlite3.db"))

	if err != nil {
		log.Fatal("Couldn't open sqlite database")
	}

	// Initialise our app-wide environment with the services/info we need.
	appCtx := &core.AppContext{
		Test: "TEST",
		DB:   db,
		//Port: os.Getenv("PORT"),
		//Host: os.Getenv("HOST"),
		// We might also have a custom log.Logger, our
		// template instance, and a config struct as fields
		// in our Env struct.
	}

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Get("/", core.AppHandleFunc(appCtx, page.Index))
	//router.Get("/page/{pageID}", page.ViewPage)

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)

	//router.Get("/auth/login", auth.Login)
	//router.Post("/auth/logout", auth.Logout)
	return router
}

func runServer(host string, port int) {
	app := CreateApp()
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Running on http://%s/ (Press CTRL+C to quit)", addr)

	http.ListenAndServe(addr, app)
}

func serverCmd() *cobra.Command {
	var srvCmd = &cobra.Command{
		Use:   "server",
		Short: "Runs the webserver",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(viper.GetString("web.host"), viper.GetInt("web.port"))
		},
	}

	srvCmd.PersistentFlags().IntP("port", "", 8080, "The port to bind to")
	srvCmd.PersistentFlags().StringP("host", "", "127.0.0.1", "The address to bind to")
	viper.BindPFlag("web.port", srvCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("web.host", srvCmd.PersistentFlags().Lookup("host"))
	return srvCmd
}
