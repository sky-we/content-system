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
	sessionId int
	rdb       *redis.Client
}

func (s *SessionAuth) Auth(ctx *gin.Context) {
	sid := ctx.GetHeader(sessionKey)
	if sid == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "用户未登录")
		return
	}
	s.connRdb()
	redisCtx := context.Background()
	_, err := s.rdb.Get(redisCtx, utils.GenAuthKey(sid)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "系统内部错误，请稍后重试")
		return
	}
	if errors.Is(err, redis.Nil) {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "用户未登录")
		return
	}
	fmt.Printf("session id = %s", utils.GenAuthKey(sid))
	ctx.Next()
}

func (s *SessionAuth) connRdb() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	s.rdb = rdb
}
