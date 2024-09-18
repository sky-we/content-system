package services

import (
	"github.com/redis/go-redis/v9"
	goflow "github.com/s8sg/goflow/v1"
	"gorm.io/gorm"
)

type CmsApp struct {
	db          *gorm.DB
	rdb         *redis.Client
	flowService *goflow.FlowService
}

func NewCmsApp(db *gorm.DB, rdb *redis.Client, flowService *goflow.FlowService) *CmsApp {
	app := &CmsApp{
		db:          db,
		rdb:         rdb,
		flowService: flowService,
	}
	return app
}
