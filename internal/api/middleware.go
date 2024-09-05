package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const sessionKey = "session_id"

type SessionAuth struct {
	sessionId int
}

func (s SessionAuth) Auth(ctx *gin.Context) {
	if ctx.GetHeader(sessionKey) == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "session id is null")
		return
	}
	fmt.Println("session id = ", s.sessionId)
	ctx.Next()
}
