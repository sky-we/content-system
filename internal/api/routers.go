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
	db := config.NewMySqlDB(config.DBConfig.MySQL)
	rdb := config.NewRdb(config.DBConfig.Redis)

	// 依赖注入
	cmsApp := services.NewCmsApp(db, rdb)

	// 鉴权中间件
	sessionMiddleware := &SessionAuth{Rdb: rdb}
	root := r.Group(rootPath).Use(sessionMiddleware.Auth)
	{
		// 服务探测
		root.GET("/cms/probe", cmsApp.Probe)

		// 内容创建
		root.POST("/cms/content/create", cmsApp.ContentCreate)

		// 内容更新
		root.POST("/cms/content/update", cmsApp.ContentUpdate)

		// 内容更新
		root.POST("/cms/content/delete", cmsApp.ContentDelete)
	}

	outRoot := r.Group(outRootPath)
	{
		// 用户注册
		outRoot.POST("/cms/register", cmsApp.Register)

		// 用户登录
		outRoot.POST("/cms/login", cmsApp.Login)
	}

}
