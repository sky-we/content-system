package main

import (
	"content-system/internal/api"
	"content-system/internal/config"
	"content-system/internal/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadDBConfig()
}

var Logger = middleware.GetLogger()

func main() {
	middleware.InitLogger()
	r := gin.Default()
	api.CmsRouters(r)
	if err := r.Run(); err != nil {
		Logger.Error("run err %v", err)
		return
	}

}
