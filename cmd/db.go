package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func main() {

}

func runMigrations() {
	fmt.Println("Running migrations...")
	dir, _ := os.Getwd()
	dbPath := filepath.Join(dir, "sqlite3.db")
	dbMigrationsPath := filepath.Join(dir, "migrations", "sqlite3")
	log.Println("Using database: ", dbPath)
	log.Println("Using migrations from: ", dbMigrationsPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Couldn't open sqlite database")
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", dbMigrationsPath),
		"someblocks.sqlite", driver)

	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Runs the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}
