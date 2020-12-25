package controllers

import (
	"net/http"
	"someblocks/internal/app"

	"github.com/go-chi/chi"
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
	c.app.HTML(w, r, "index", app.D{"hello": "json"})
}

func (c *PageController) PageView(w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "page", gin.H{
	//	"Title": "Hello Page!",
	//	"Body": "Mah body is dat",
	//})
	pageID := chi.URLParam(r, "pageID")
	c.app.HTML(w, r, "page/page", app.D{"pageID": pageID})
}
