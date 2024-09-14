package services

import (
	"content-system/internal/dao"
	"content-system/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContentDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

type ContentDeleteRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func (app *CmsApp) ContentDelete(ctx *gin.Context) {
	var contentDeleteReq ContentDeleteReq
	if err := ctx.ShouldBindJSON(&contentDeleteReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "参数错误", "error": err.Error()})
		return
	}
	contentDetailDao := dao.NewContentDetailDao(app.db)
	exists, err := contentDetailDao.IsExist(contentDeleteReq.ID)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf("[ID=%d]内容不存在]", contentDeleteReq.ID)})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}

	if err := contentDetailDao.Delete(contentDeleteReq.ID, &model.ContentDetail{}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ContentDeleteRsp{
		Code: 0,
		Msg:  "success",
		Data: fmt.Sprintf("ID %d delete", contentDeleteReq.ID),
	})

}
