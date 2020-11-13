package core

import (
	"github.com/jmoiron/sqlx"
)


// This holds the app context
type AppContext struct {
	DB   *sqlx.DB
	Port string
	Host string
}
