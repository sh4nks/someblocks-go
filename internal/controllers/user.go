package controllers

import (
	"net/http"
	"someblocks/internal/app"
	"someblocks/internal/forms"
	"someblocks/internal/middleware"
	"someblocks/internal/models"

	"github.com/go-chi/chi"
)

func NewUserController(app *app.App, userService *models.UserService) *UserController {
	return &UserController{
		app:         app,
		userService: userService,
	}
}

type UserController struct {
	app         *app.App
	userService *models.UserService
}

func (c *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/profile", app.D{})
}

//
// User Settings "Profile"
//

func (c *UserController) SettingsProfile(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_profile", app.D{})
}

func (c *UserController) SettingsProfileUpdate(w http.ResponseWriter, r *http.Request) {
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

	c.userService.Update()

	c.app.Flash(r, "Profile updated!", "success")
	http.Redirect(w, r, "/user/settings/profile", http.StatusSeeOther)
}

//
// User Settings "Account"
//

func (c *UserController) SettingsAccount(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_account", app.D{})
}

func (c *UserController) SettingsAccountUpdate(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_account", app.D{})
}

func (c *UserController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_profile", app.D{})
}

//
// User Settings "Security"
//

func (c *UserController) SettingsSecurity(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_security", app.D{})
}

func (c *UserController) SettingsSecurityPost(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_security", app.D{})
}

//
// User Settings "Notification"
//

func (c *UserController) SettingsNotification(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "user/settings_notification", app.D{})
}

func (c *UserController) Routes(userMw *middleware.User) chi.Router {
	router := chi.NewRouter()
	router.Use(userMw.LoginRequired)

	// Profile
	router.Get("/profile", c.Profile)

	// Settings Profile
	router.Get("/settings", c.SettingsProfile)
	router.Get("/settings/profile", c.SettingsProfile)
	router.Post("/settings/profile", c.SettingsProfileUpdate)

	// Settings Account
	router.Get("/settings/account", c.SettingsAccount)
	router.Post("/settings/account", c.SettingsAccountUpdate)
	router.Post("/settings/account/delete", c.DeleteAccount)

	router.Get("/settings/security", c.SettingsSecurity)
	router.Get("/settings/notification", c.SettingsNotification)
	return router
}
