package auth

type LoginForm struct {
	Login      string `form:"login" binding:"required"`
	Password   string `form:"password" binding:"required"`
	RememberMe string    `form:"rememberMe"`
}

type RegisterForm struct {
	username        string
	email           string
	password        string
	confirmPassword string
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
