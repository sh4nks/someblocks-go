package app

import (
	"github.com/jmoiron/sqlx"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type App struct {
	DB     *sqlx.DB
}

func New() *App {
	return &App{}
}
