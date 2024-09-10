package api

import (
	"content-system/internal/config"
	"content-system/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	rootPath    = "/api"
	outRootPath = "/out/api"
)

func CmsRouters(r *gin.Engine) {

	db := config.NewMySqlDB(config.DBConfig.MySQL)
	rdb := config.NewRdb(config.DBConfig.Redis)

	cmsApp := services.NewCmsApp(db, rdb)

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
	//defer func() {
	//	CloseDB(db, rdb)
	//}()

}

func CloseDB(db *gorm.DB, rdb *redis.Client) {

	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		fmt.Println("close mysql db error")
		panic(err)
	}
	if err := rdb.Close(); err != nil {
		fmt.Println("close redis db error")
		panic(err)
	}

}
