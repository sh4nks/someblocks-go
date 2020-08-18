package page

import (
	"github.com/gin-gonic/gin"
)


func IndexView(ctx *gin.Context) {
	ctx.HTML(200, "index", gin.H{})
}

func PageView(ctx *gin.Context) {
	ctx.HTML(200, "page", gin.H{
		"Title": "Hello Page!",
		"Body": "Mah body is dat",
	})
}
