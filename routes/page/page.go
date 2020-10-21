package page

import (
	"github.com/go-chi/chi"
	"net/http"
	"someblocks/core"
)

func Index(appCtx *core.AppContext, w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "index", gin.H{})
	w.Write([]byte(appCtx.Test))
}

func ViewPage(w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "page", gin.H{
	//	"Title": "Hello Page!",
	//	"Body": "Mah body is dat",
	//})
	pageID := chi.URLParam(r, "pageID")
	w.Write([]byte(pageID))
}
