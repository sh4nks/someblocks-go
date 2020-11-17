package app

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"someblocks/internal/core"

	"someblocks/internal/routes/auth"
	"someblocks/internal/routes/page"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

type App struct {
	Ctx    *core.AppContext
	Routes http.Handler
}

func (app *App) CreateApp() http.Handler {

	// Create the AppContext if it doesn't exist
	if app.Ctx == nil {
		app.Ctx = &core.AppContext{}

		db, err := app.SetupDatabase(viper.GetString("database.driver"))

		if err != nil {
			log.Fatal(err)
		}

		app.Ctx.DB = db
	}

	// Setup "Template Engine" AKA renderer
	render := render.New(render.Options{
				Layout: "layout",
				Extensions: []string{".html"},
			})
	app.Ctx.UseRender(render)

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Get("/", core.AppHandleFunc(app.Ctx, page.Index))
	router.Get("/page/{pageID}", core.AppHandleFunc(app.Ctx, page.ViewPage))

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)

	router.Get("/auth/login", core.AppHandleFunc(app.Ctx, auth.Login))
	router.Post("/auth/logout", core.AppHandleFunc(app.Ctx, auth.Logout))

	app.Routes = router
	return router
}

func (app *App) SetupDatabase(drivername string) (*sqlx.DB, error) {
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

func (app *App) Migrate() {
	log.Println("Running migrations...")
	db, err := app.SetupDatabase(viper.GetString("database.driver"))
	if err != nil {
		log.Fatal(err)
	}

	drivername := viper.GetString("database.driver")
	dir := filepath.Join(core.GetAppDir(), "migrations")
	if drivername == "sqlite" || drivername == "sqlite3" {
		migrationsPath := fmt.Sprintf("file:///%s", filepath.Join(dir, "sqlite3"))

		log.Println("Using migrations from: ", migrationsPath)

		driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite", driver)

		if err != nil {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}

	} else if drivername == "postgres" {
		migrationsPath := fmt.Sprintf("file:///%s", filepath.Join(dir, "postgres"))
		log.Println("Using migrations from: ", migrationsPath)

		driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)

		if err != nil {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}
}
