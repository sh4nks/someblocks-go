package core

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

// This holds the app context
type AppContext struct {
	Render *render.Render
	DB     *sqlx.DB
	Port   string
	Host   string
}


// H is a shortcut for map[string]interface{}
type H map[string]interface{}

func (c *AppContext) JSON(w http.ResponseWriter, status int, v interface{}) {
	c.Render.JSON(w, status, v)
}

func (c *AppContext) Text(w http.ResponseWriter, status int, v string) {
	c.Render.Text(w, status, v)
}

func (c *AppContext) HTML(w http.ResponseWriter, status int, tmpl string, data interface{}) {
	c.Render.HTML(w, status, tmpl, data)
}
