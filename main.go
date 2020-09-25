package main

import (
	//"github.com/foolin/goview"
	//"github.com/gin-gonic/gin"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"someblocks/core"
	//"someblocks/controllers/blog"
	"someblocks/controllers/page"
)



func Routes() http.Handler {

    // Initialise our app-wide environment with the services/info we need.
    appCtx := &handler.AppContext{
    	Test: "TEST",
        //DB: db,
        //Port: os.Getenv("PORT"),
        //Host: os.Getenv("HOST"),
        // We might also have a custom log.Logger, our
        // template instance, and a config struct as fields
        // in our Env struct.
    }


	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Get("/", handler.AppHandleFunc(appCtx, page.Index))
	//router.Get("/page/{pageID}", page.ViewPage)

	//router.Get("/blog", blog.Index)
	//router.Get("/blog/{blogID}", blog.ViewPost)

	//router.Get("/auth/login", auth.Login)
	//router.Post("/auth/logout", auth.Logout)
	return router
}



func main() {
	router := Routes()
	//router.Run("127.0.0.1:8080")
	http.ListenAndServe(":3333", router)
}
