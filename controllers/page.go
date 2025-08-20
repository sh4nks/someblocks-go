package controllers

import (
	"net/http"
	"someblocks/app"
	"someblocks/context"

	"github.com/go-chi/chi/v5"
)

func NewPageController(app *app.App) *PageController {
	return &PageController{
		app: app,
	}
}

type PageController struct {
	app *app.App
}

func (c *PageController) PageIndex(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "index", app.D{})
}

func (c *PageController) PageView(w http.ResponseWriter, r *http.Request) {
	user := context.CurrentUser(r.Context())
	pageID := chi.URLParam(r, "pageID")
	c.app.HTML(w, r, "page/page", app.D{"pageID": pageID, "userFromContext": user})
}
