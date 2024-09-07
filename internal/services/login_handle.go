package services

import (
	"content-system/internal/dao"
	"content-system/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginReq struct {
	UserId   string `json:"user_id" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
}
type LoginRsp struct {
	SessionId string `json:"session_id" binding:"required"`
	UserId    string `json:"user_id" binding:"required"`
	NickName  string `json:"nick_name" binding:"required"`
}

func (app *CmsApp) Login(ctx *gin.Context) {

	var loginReq LoginReq

	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountDao := dao.NewAccountDao(app.db)
	account, err := accountDao.FindByUserId(loginReq.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "请输入正确的用户ID")
		return

	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginReq.PassWord)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "请输入正确的密码")
		return

	}
	sessionId, err := app.genSessionId(ctx, account.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "系统内部错误，请稍后重试")

	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0,
		"msg": "login ok",
		"data": &LoginRsp{
			SessionId: sessionId,
			UserId:    account.UserId,
			NickName:  account.NickName,
		}})

}

func (app *CmsApp) genSessionId(context context.Context, userId string) (string, error) {
	sessionId := uuid.New().String()
	//sessionKey := utils.GenSessionKey(sessionId)
	//if err := app.rdb.Set(context, utils.SessionKey, sessionId, time.Hour*8).Err(); err != nil {
	//	return "", err
	//}
	//fmt.Println("sessionId", sessionId)
	authKey := utils.GenAuthKey(sessionId)
	if err := app.rdb.Set(context, authKey, time.Now().Unix(), time.Hour*8).Err(); err != nil {
		return "", err
	}
	return sessionId, nil
}
