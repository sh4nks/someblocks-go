package app

import (
	"fmt"

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
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

type App struct {
	Ctx    *core.AppContext
	Routes http.Handler
}

func NewApp() *App {
	app := App{}
	app.InitApp()
	return &app
}

func (app *App) InitApp() {

	db, err := app.SetupDatabase(viper.GetString("database.driver"))

	if err != nil {
		log.Fatal().Err(err).Msg("An error occured while seting up the database")
	}

	// Setup "Template Engine" AKA renderer
	render := render.New(render.Options{
				Layout: "layout",
				Extensions: []string{".html"},
			})

	if app.Ctx == nil {
		app.Ctx = &core.AppContext{
			DB: db,
			Render: render,
		}
	}

	app.autoMigrate()
	app.registerRoutes()
}

func (app *App) registerRoutes() {
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
}

func (app *App) SetupDatabase(drivername string) (*sqlx.DB, error) {
	if drivername == "sqlite3" || drivername == "sqlite" {
		dir := core.GetAppDir()
		connStr := filepath.Join(dir, viper.GetString("database.dbname"))

		log.Debug().Msgf("Using sqlite3 with following connection string: %s", connStr)

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
		log.Debug().Msgf("Using postgres with following connection string: %s", redactedStr)

		db, err := sqlx.Connect("postgres", connStr)

		if err != nil {
			return nil, fmt.Errorf("Couldn't connect to postgres database: %s", err)
		}

		return db, nil
	}

	return nil, fmt.Errorf("%s is not supported", drivername)
}

func (app *App) autoMigrate() {
	log.Info().Msg("Running auto migrations...")
	mustAppContext(app)
	drivername := viper.GetString("database.driver")
	if drivername == "sqlite" || drivername == "sqlite3" {
		migrateSQLite(app.Ctx.DB)
	} else if drivername == "postgres" {
		migratePostgres(app.Ctx.DB)
	}
}

func (app *App) Migrate() {
	log.Info().Msg("Running migrations...")

	db, err := app.SetupDatabase(viper.GetString("database.driver"))
	if err != nil {
		log.Fatal().Err(err).Msg("An error occured during the database setup")
	}

	drivername := viper.GetString("database.driver")
	if drivername == "sqlite" || drivername == "sqlite3" {
		migrateSQLite(db)
	} else if drivername == "postgres" {
		migratePostgres(db)
	}
}


func migratePostgres(db *sqlx.DB) {
	dir := filepath.Join(core.GetAppDir(), "migrations")
	migrationsPath := fmt.Sprintf("file:///%s", filepath.Join(dir, "postgres"))
	log.Info().Msgf("Using migrations from: %s", migrationsPath)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("An error occured while trying to use postgres")
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)

	if err != nil {
		log.Fatal().Err(err).Msg("An error occured while creating a new migrate instance")
	}

	err = m.Up()
	if (err != nil && err.Error() == "no change") {
		log.Info().Msg("No changes")
	} else if err != nil {
		log.Fatal().Err(err).Msg("An error occured while running the migrations")
	} else {
		log.Info().Msg("Database schema updated")
	}
}


func migrateSQLite(db *sqlx.DB) {
	dir := filepath.Join(core.GetAppDir(), "migrations")
	migrationsPath := fmt.Sprintf("file:///%s", filepath.Join(dir, "sqlite3"))

	log.Info().Msgf("Using migrations from: %s", migrationsPath)

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("An error occured while trying to use sqlite")
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("An error occured while creating a new migrate instance")
	}

	err = m.Up()
	if (err != nil && err.Error() == "no change") {
		log.Info().Msg("No changes")
	} else if err != nil {
		log.Fatal().Err(err).Msg("An error occured while running the migrations")
	} else {
		log.Info().Msg("Database schema updated")
	}
}


func mustAppContext(app *App) {
	if app.Ctx == nil {
		log.Fatal().Msg("Not running inside an app context.")
	}
}
