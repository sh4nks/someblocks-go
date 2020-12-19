package controllers

import (
	"someblocks/internal/app"

	"github.com/labstack/echo/v4"
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

func (auth *AuthController) Login(c echo.Context) error {
	return c.Render(200, "auth/login", Data{"Title": "Login"})
}

func (auth *AuthController) LoginPost(c echo.Context) error {
	//var login LoginForm
	return c.String(200, "POST LOGIN")
}

func (auth *AuthController) Logout(c echo.Context) error {
	return c.String(200, "Hello Logout")
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
