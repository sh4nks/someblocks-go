package blog

import (
	"net/http"
	//"github.com/go-chi/chi"
	"someblocks/core"
)

func Index(appCtx *handler.AppContext, w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "index", gin.H{})
	w.Write([]byte("Hello World!"));
}
