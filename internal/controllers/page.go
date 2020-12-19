package controllers

import (
	"someblocks/internal/app"

	"github.com/labstack/echo/v4"
)

func NewPageController(app *app.App) *PageController {
	return &PageController{
		app: app,
	}
}

type PageController struct {
	app *app.App
}

func (page *PageController) PageIndex(c echo.Context) error {
	return c.Render(200, "index", Data{"hello": "json"})
}

func (page *PageController) PageView(c echo.Context) error {
	//ctx.HTML(200, "page", gin.H{
	//	"Title": "Hello Page!",
	//	"Body": "Mah body is dat",
	//})
	pageID := c.Param("pageID")
	return c.Render(200, "page/page", Data{"pageID": pageID})
}
