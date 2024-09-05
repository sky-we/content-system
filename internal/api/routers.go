package api

import (
	"content-system/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	rootPath    = "/api"
	outRootPath = "/out/api"
)

func CmsRouters(r *gin.Engine) {
	cmsApp := services.NewCmsApp()
	sessionMiddleware := &SessionAuth{}
	root := r.Group(rootPath).Use(sessionMiddleware.Auth)
	{
		root.GET("/cms/probe", cmsApp.Probe)

	}

	outRoot := r.Group(outRootPath)
	{
		outRoot.POST("/cms/register", cmsApp.Register)
	}

}
