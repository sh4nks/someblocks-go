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
		c.app.Flash(r, "Login error", "danger")
		c.app.HTML(w, r, "auth/login", app.D{
			"Title": "Login",
			"Form":  form,
		})
		return
	}

	id, err := c.userService.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidLoginCredentials {
		c.app.Flash(r, "Login error", "danger")
		form.Errors.Add("generic", "Email or Password is incorrect")
		c.app.HTML(w, r, "auth/login", app.D{
			"Title": "Login",
			"Form":  form,
		})
		return
	} else if err != nil {
		c.app.ServerError(w, err)
		return
	}

	if form.Get("rememberMe") == "on" {
		c.app.Session.RememberMe(r.Context(), true)
	}

	//c.app.Flash(r, "Logged in!", "success")
	c.app.Session.Put(r.Context(), "userId", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) LogoutPost(w http.ResponseWriter, r *http.Request) {
	c.app.Flash(r, "You've been logged out successfully!", "success")
	c.app.Session.Remove(r.Context(), "userId")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	if !form.Valid(c.userService) {
		c.app.Flash(r, "Something went wrong", "danger")
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

	c.app.Flash(r, "Registered!", "success")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
