package api

import (
	"content-system/internal/config"
	"content-system/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	rootPath    = "/api"
	outRootPath = "/out/api"
)

func CmsRouters(r *gin.Engine) {
	db := config.dbConfig
	cmsApp := services.NewCmsApp()
	sessionMiddleware := &SessionAuth{}
	root := r.Group(rootPath).Use(sessionMiddleware.Auth)
	{
		// 服务探测
		root.GET("/cms/probe", cmsApp.Probe)
	}

	outRoot := r.Group(outRootPath)
	{
		// 用户注册
		outRoot.POST("/cms/register", cmsApp.Register)

		// 用户登录
		outRoot.POST("/cms/login", cmsApp.Login)
	}

}
