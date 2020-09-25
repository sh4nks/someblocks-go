package page

import (
	"net/http"
	"github.com/go-chi/chi"
	"someblocks/core"
)


func Index(appCtx *handler.AppContext, w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "index", gin.H{})
	w.Write([]byte(appCtx.Test));
}

func ViewPage(w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "page", gin.H{
	//	"Title": "Hello Page!",
	//	"Body": "Mah body is dat",
	//})
	pageID := chi.URLParam(r, "pageID")
	w.Write([]byte(pageID));
}
