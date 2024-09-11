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

	if app.IsLogin(context.Background(), loginReq.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "用户已登录"})
		return
	}

	accountDao := dao.NewAccountDao(app.db)
	account, err := accountDao.FindByUserId(loginReq.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "请输入正确的用户ID"})
		return

	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginReq.PassWord)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "用户密码错误"})
		return

	}
	sessionId, err := app.genSessionId(context.Background(), account.UserId)
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

func (app *CmsApp) genSessionId(ctx context.Context, userId string) (string, error) {
	sessionId := uuid.New().String()
	if err := app.rdb.Set(ctx, utils.GenSessionKey(userId), sessionId, time.Hour*8).Err(); err != nil {
		return "", err
	}
	if err := app.rdb.Set(ctx, utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err(); err != nil {
		return "", err
	}
	return sessionId, nil
}

func (app *CmsApp) IsLogin(ctx context.Context, userId string) bool {
	exists, err := app.rdb.Exists(ctx, utils.GenSessionKey(userId)).Result()
	if err != nil {
		panic(err)
	}
	return exists > 0

}
