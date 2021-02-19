package database

import (
    "fmt"
    "net/url"

    "path/filepath"
    "someblocks/internal/config"
    "someblocks/internal/models"
    "someblocks/pkg/utils"

    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "github.com/rs/zerolog/log"
)

//  case drivers.DatabaseDriverSqlite:
//      sqliteAddress, err := getSqliteAddress()
//      if err != nil {
//          return nil, err
//      }
//      log.Printf("Opening SQLITE database: %s", sqliteAddress)
//      databaseDialect = sqlite.Open(sqliteAddress.String())//

//  case drivers.DatabaseDriverPostgres:
//      postgresAddress, err := getPostgresAddress()
//      if err != nil {
//          return nil, err
//      }
//      log.Printf("Connecting to POSTGRES database: %s", postgresAddress.Redacted())
//      databaseDialect = postgres.Open(postgresAddress.String())
//  }

func getSqliteAddress(cfg *config.Config) (*url.URL, error) {

    if cfg.Database.URL == "" {
        dir := utils.GetExecDir()
        cfg.Database.URL = filepath.Join(dir, "someblocks.sqlite")
    }

    address, err := url.Parse(cfg.Database.URL)
    if err != nil {
        return nil, fmt.Errorf("Could not parse sqlite url (%s): %s", cfg.Database.URL, err)
    }

    return address, nil
}

func getPostgresAddress(cfg *config.Config) (*url.URL, error) {
    if cfg.Database.URL == "" {
        return nil, fmt.Errorf("No connection URL provided for driver postgres")
    }

    address, err := url.Parse(cfg.Database.URL)
    if err != nil {
        return nil, fmt.Errorf("Could not parse postgres url: %s", err)
    }

    return address, nil
}

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {

    var dialect gorm.Dialector
    switch driver := cfg.Database.Driver; driver {
    case "sqlite":
        connStr, err := getSqliteAddress(cfg)
        if err != nil {
            return nil, err
        }

        log.Info().Msgf("Connecting to sqlite database: %s", connStr)
        dialect = sqlite.Open(connStr.String())
    case "postgres":
        connStr, err := getPostgresAddress(cfg)
        if err != nil {
            return nil, err
        }

        log.Info().Msg("Connecting to postgres database")
        dialect = postgres.Open(connStr.String())
    default:
        return nil, fmt.Errorf("%s is not supported", driver)
    }

    db, err := gorm.Open(dialect, &gorm.Config{})
    return db, err
}

func Migrate(db *gorm.DB) error {
    log.Info().Msg("Running database migrations...")
    db.AutoMigrate(
        &models.User{},
    )
    return nil
}

func SetupAndMigrate(cfg *config.Config) (*gorm.DB, error) {
    db, err := SetupDatabase(cfg)
    if err != nil {
        return nil, err
    }

    err = Migrate(db)
    if err != nil {
        return nil, err
    }

    return db, nil
}
