package core

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

// This holds the app context
type AppContext struct {
	render *render.Render
	DB     *sqlx.DB
	Port   string
	Host   string
}


// H is a shortcut for map[string]interface{}
type H map[string]interface{}

func (c *AppContext) UseRender(render *render.Render) {
	c.render = render
}

func (c *AppContext) JSON(w http.ResponseWriter, status int, v interface{}) {
	c.render.JSON(w, status, v)
}

func (c *AppContext) Text(w http.ResponseWriter, status int, v string) {
	c.render.Text(w, status, v)
}

func (c *AppContext) HTML(w http.ResponseWriter, status int, tmpl string, data interface{}) {
	c.render.HTML(w, status, tmpl, data)
}
