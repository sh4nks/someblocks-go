package controllers

import (
	"net/http"
	"someblocks/internal/app"
	"someblocks/internal/forms"
	"someblocks/internal/models"
)

func NewAuthController(app *app.App, userService *models.UserService) *AuthController {
	return &AuthController{
		app:         app,
		userService: userService,
	}
}

type AuthController struct {
	app         *app.App
	userService *models.UserService
}

type LoginForm struct {
	Login      string `form:"login" binding:"required"`
	Password   string `form:"password" binding:"required"`
	RememberMe string `form:"rememberMe"`
}

type RegisterForm struct {
	username        string
	email           string
	password        string
	confirmPassword string
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "auth/login", app.D{
		"Title": "Login",
	})
}

func (c *AuthController) LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		c.app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewLoginForm(r.PostForm)
	if !form.Valid() {
		c.app.Session.SetFlash(r.Context(), "danger", "Login error")
		c.app.HTML(w, r, "auth/login", app.D{
			"Title": "Login",
			"Form":  form,
		})
		return
	}
	c.userService.Authenticate(form.Get("email"), form.Get("password"))

	c.app.Session.SetFlash(r.Context(), "success", "Logged in!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) LogoutPost(w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "index", gin.H{})
	w.Write([]byte("Hello Logout"))
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	c.app.HTML(w, r, "auth/register", app.D{
		"Title": "Register",
	})
}

func (c *AuthController) RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		c.app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewRegisterForm(r.PostForm)
	if !form.Valid() {
		c.app.Session.SetFlash(r.Context(), "danger", "Something went wrong")
		c.app.HTML(w, r, "auth/register", app.D{
			"Title": "Register",
			"Form":  form,
		})
		return
	}

	c.userService.Insert(
		form.Get("username"),
		form.Get("email"),
		form.Get("password"),
	)

	c.app.Session.SetFlash(r.Context(), "success", "Registered!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
