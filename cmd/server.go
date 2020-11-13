package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"someblocks/core"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func SetupDb(drivername string) (*sqlx.DB, error) {
	if drivername == "sqlite3" || drivername == "sqlite" {
		dir := core.GetAppDir()
		connStr := filepath.Join(dir, viper.GetString("database.dbname"))

		log.Printf("Using sqlite3 with following connection string: %s", connStr)

		db, err := sqlx.Open("sqlite3", connStr)

		if err != nil {
			return nil, fmt.Errorf("Couldn't open sqlite3 database: %s", err)
		}

		return db, nil
	} else if drivername == "postgres" {
		password := ""
		if viper.GetString("database.password") != "" {
			password = fmt.Sprintf("password=%s", viper.GetString("database.password"))
		}

		host := viper.GetString("database.host")
		port := viper.GetInt("database.port")
		dbname := viper.GetString("database.dbname")
		user := viper.GetString("database.username")

		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s %s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)

		var redactedStr string
		if password != "" {
			redactedStr = fmt.Sprintf(
				"host=%s port=%d user=%s password=***** dbname=%s sslmode=disable",
				host, port, user, dbname,
			)
		} else {
			redactedStr = connStr
		}
		log.Printf("Using postgres with following connection string: %s", redactedStr)

		db, err := sqlx.Connect("postgres", connStr)

		if err != nil {
			return nil, fmt.Errorf("Couldn't connect to postgres database: %s", err)
		}

		return db, nil
	}

	return nil, fmt.Errorf("%s is not supported", drivername)
}

func CreateApp() http.Handler {

	db, err := SetupDb(viper.GetString("database.driver"))

	if err != nil {
		log.Fatal(err)
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
