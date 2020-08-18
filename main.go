package main

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"peterjustin.com/website/modules/auth"
	"peterjustin.com/website/modules/page"
)



func Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.HTMLRender = ginview.New(goview.Config{
		Root: "templates",
		Extension: ".html",
		Master: "layout",
		// if cache disabled, auto reload template file for debug.
		DisableCache: true,
	})

	router.Static("/static", "./assets")

	pageRouter := router.Group("/")
	pageRouter.GET("/", page.IndexView)
	pageRouter.GET("/p/:pageName", page.PageView)

	authRouter := router.Group("/auth")
	authRouter.POST("/login", auth.LoginPost)
	authRouter.GET("/login", auth.Login)
	authRouter.POST("/logout")
	authRouter.GET("/register")
	authRouter.POST("/register")

	blogRouter := router.Group("/blog")
	blogRouter.GET("/")



	//router := chi.NewRouter()
	//router.Use(
	//	middleware.RequestID,
	//	middleware.Logger,
	//	middleware.RedirectSlashes,
	//	middleware.Recoverer,
	//)

	//router.Route("/", func(r chi.Router) {
	//	r.Mount("/", page.Routes())
	//	r.Mount("/auth", auth.Routes())
	//	r.Mount("/admin", admin.Routes())
	//	r.Mount("/blog", blog.Routes())
	//})

	return router
}



func main() {
	router := Routes()
	router.Run("127.0.0.1:8080")
}
