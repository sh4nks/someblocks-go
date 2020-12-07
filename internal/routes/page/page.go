package page

import (
	"net/http"
	"someblocks/internal/core"

	"github.com/go-chi/chi"
)

func Index(ctx *core.AppContext, w http.ResponseWriter, r *http.Request) {
	ctx.HTML(w, r, 200, "index", core.H{"hello": "json"})
}

func ViewPage(ctx *core.AppContext, w http.ResponseWriter, r *http.Request) {
	//ctx.HTML(200, "page", gin.H{
	//	"Title": "Hello Page!",
	//	"Body": "Mah body is dat",
	//})
	pageID := chi.URLParam(r, "pageID")
	ctx.HTML(w, r, 200, "page/page", core.H{"pageID": pageID})
}
