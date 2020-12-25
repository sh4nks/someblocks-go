package controllers

import (
	"net/http"
	"someblocks/internal/app"
	"someblocks/internal/forms"
)

func NewAuthController(app *app.App) *AuthController {
	return &AuthController{
		app: app,
	}
}

type AuthController struct {
	app *app.App
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

	//form.Get("Email")
	//form.Get("Password")

	c.app.Session.SetFlash(r.Context(), "success", "Logged in!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) LogoutPost(w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "index", gin.H{})
	w.Write([]byte("Hello Logout"))
}

/*
func Login(ctx *gin.Context) {
	ctx.HTML(200, "login", gin.H{
		"title": "Login",
	})
}

func LoginPost(ctx *gin.Context) {
	var login LoginForm

	if err := ctx.ShouldBind(&login); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("ELSE")
		log.Println(login.Login)
		log.Println(login.Password)
		log.Println(login.RememberMe)
	}
}
*/
