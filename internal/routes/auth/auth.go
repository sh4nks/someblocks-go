package auth

import (
	"net/http"
	"someblocks/internal/core"

	"github.com/go-chi/chi"
)

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

func Login(ctx *core.AppContext, w http.ResponseWriter, r *http.Request) {
	ctx.HTML(w, 200, "login", core.H{"Title": "Login"})
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	pageID := chi.URLParam(r, "pageID")
	w.Write([]byte(pageID))
}

func Logout(ctx *core.AppContext, w http.ResponseWriter, r *http.Request) {
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
