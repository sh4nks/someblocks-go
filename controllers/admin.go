package controllers

import (
	"net/http"
	"someblocks/app"
	"someblocks/forms"
	"someblocks/middleware"
	"someblocks/models"

	"github.com/go-chi/chi/v5"
)

func NewAdminController(app *app.App, userService *models.UserService) *AdminController {
	return &AdminController{
		app:         app,
		userService: userService,
	}
}

type AdminController struct {
	app         *app.App
	userService *models.UserService
}

//
// User Settings "Profile"
//

func (c *AdminController) Dashboard(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "admin/dashboard", app.D{})
}

func (c *AdminController) SettingsProfileUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		c.app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewProfileForm(r.PostForm)
	if !form.Valid(c.userService) {
		c.app.Flash(r, "Something went wrong", "danger")
		c.app.HTML(w, r, "user/settings_profile", app.D{
			"Form": form,
		})
		return
	}
	//c.userService.Update()

	c.app.Flash(r, "Profile updated!", "success")
	http.Redirect(w, r, "/user/settings/profile", http.StatusSeeOther)
}

func (c *AdminController) Routes(userMw *middleware.User) chi.Router {
	router := chi.NewRouter()
	router.Use(userMw.AdminRequired)

	// Settings Profile
	router.Get("/admin", c.Dashboard)
	router.Post("/settings/profile", c.SettingsProfileUpdate)
	return router
}
