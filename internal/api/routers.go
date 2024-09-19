package api

import (
	"content-system/internal/config"
	"content-system/internal/middleware"
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
	flowService := config.NewFlowService(config.DBConfig.FlowService, db)

	// 依赖注入
	cmsApp := services.NewCmsApp(db, rdb, flowService)

	// 启动内容加工Flow
	cmsApp.StartFlow(flowService)

	// 鉴权中间件
	sessionMiddleware := &middleware.SessionAuth{Rdb: rdb}

	// 路由
	router := r.Group(rootPath)
	router.Use(sessionMiddleware.Auth)
	{
		// 服务探测
		router.GET("/cms/probe", cmsApp.Probe)

		// 内容创建
		router.POST("/cms/content/create", cmsApp.ContentCreate)

		// 内容更新
		router.POST("/cms/content/update", cmsApp.ContentUpdate)

		// 内容删除
		router.POST("/cms/content/delete", cmsApp.ContentDelete)

		// 内容查询
		router.POST("/cms/content/find", cmsApp.ContentFind)
	}

	outRoot := r.Group(outRootPath)
	{
		// 用户注册
		outRoot.POST("/cms/register", cmsApp.Register)

		// 用户登录
		outRoot.POST("/cms/login", cmsApp.Login)
	}

}
