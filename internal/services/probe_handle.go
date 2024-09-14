package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 模型绑定

func (app *CmsApp) Probe(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{"Message": "Service Online!"},
	})

}
