package models

import (
	"fmt"

	"path/filepath"
	"someblocks/internal/config"
	"someblocks/pkg/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func SetupDatabase(cfg *config.Config) (*sqlx.DB, error) {

	switch driver := cfg.Database.Driver; driver {
	case "sqlite3", "sqlite":
		dir := utils.GetExecDir()
		connStr := filepath.Join(dir, cfg.Database.Dbname)

		log.Debug().Msgf("Using sqlite3 with following connection string: %s", connStr)

		db, err := sqlx.Open("sqlite3", connStr)

		if err != nil {
			return nil, fmt.Errorf("Couldn't open sqlite3 database: %s", err)
		}
		return db, nil

	case "postgres":
		password := ""
		if cfg.Database.Password != "" {
			password = fmt.Sprintf("password=%s", cfg.Database.Password)
		}

		host := cfg.Database.Host
		port := cfg.Database.Port
		dbname := cfg.Database.Dbname
		user := cfg.Database.Username

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

	default:
		return nil, fmt.Errorf("%s is not supported", driver)
	}

}

func Migrate(db *sqlx.DB, cfg *config.Config) {
	log.Info().Msg("Running migrations...")

	switch driver := cfg.Database.Driver; driver {
	case "sqlite3", "sqlite":
		migrateSQLite(db)
	case "postgres":
		migratePostgres(db)
	}
}

func SetupAndMigrate(cfg *config.Config) *sqlx.DB {
	db, err := SetupDatabase(cfg)
	if err != nil {
		log.Fatal().Msgf("An error occured while seting up the database: %s", err)
	}
	Migrate(db, cfg)
	return db
}

func migratePostgres(db *sqlx.DB) {
	dir := filepath.Join(utils.GetExecDir(), "migrations")
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
	if err != nil && err.Error() == "no change" {
		log.Info().Msg("Migrations are up to date")
	} else if err != nil {
		log.Fatal().Err(err).Msg("An error occured while running the migrations")
	} else {
		log.Info().Msg("Database schema updated")
	}
}

func migrateSQLite(db *sqlx.DB) {
	dir := filepath.Join(utils.GetExecDir(), "migrations")
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
	if err != nil && err.Error() == "no change" {
		log.Info().Msg("Migrations are up to date")
	} else if err != nil {
		log.Fatal().Err(err).Msg("An error occured while running the migrations")
	} else {
		log.Info().Msg("Database schema updated")
	}
}
