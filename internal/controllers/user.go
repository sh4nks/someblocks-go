package controllers

import (
	"net/http"
	"someblocks/internal/app"
)

func NewUserController(app *app.App) *UserController {
	return &UserController{
		app: app,
	}
}

type UserController struct {
	app *app.App
}

func (c *UserController) UserProfile(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/profile", app.D{})
}

func (c *UserController) UserSettings(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings", app.D{})
}
