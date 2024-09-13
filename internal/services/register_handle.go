package services

import (
	"content-system/internal/dao"
	"content-system/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type RegisterReq struct {
	UserId   string `json:"user_id" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
}

type RegisterRsp struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"msg" binding:"required"`
	Data    string `json:"data" binding:"required"`
}

func (app *CmsApp) Register(ctx *gin.Context) {
	var registerReq RegisterReq
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// 密码加密
	hashedPassword, err := encryptPassword(registerReq.PassWord)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
	}

	// 判断账号是否存在
	accountDao := dao.NewAccountDao(app.db)
	exist, err := accountDao.IsExist(registerReq.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		fmt.Println("error", err.Error())
		return
	}
	if exist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "账号已存在"})
		return
	}

	// do create
	nowTime := time.Now()
	if err := accountDao.Create(model.Account{
		UserId:   registerReq.UserId,
		Password: hashedPassword,
		NickName: registerReq.NickName,
		Ct:       nowTime,
		Ut:       nowTime,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
	}

	// response ok
	ctx.JSON(http.StatusOK,
		&RegisterRsp{
			Code:    0,
			Message: "ok",
			Data:    fmt.Sprintf("%s register ok", registerReq.UserId),
		})

}
func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("register user password encrypt failed,error=[%v]", err)
		return "", err
	}
	return string(hashedPassword), nil

}
