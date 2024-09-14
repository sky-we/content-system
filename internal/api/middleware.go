package api

import (
	"content-system/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

const sessionKey = "session_id"

type SessionAuth struct {
	SessionId int
	Rdb       *redis.Client
}

func (s *SessionAuth) Auth(ctx *gin.Context) {
	sid := ctx.GetHeader(sessionKey)
	if sid == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "用户未登录"})
		return
	}
	redisCtx := context.Background()
	_, err := s.Rdb.Get(redisCtx, utils.GenAuthKey(sid)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "服务器内部错误", "err": err.Error()})
		return
	}
	if errors.Is(err, redis.Nil) {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "用户未登录"})
		return
	}
	fmt.Printf("session id = %s", utils.GenAuthKey(sid))
	ctx.Next()
}
