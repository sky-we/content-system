package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CmsApp struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCmsApp(db *gorm.DB, rdb *redis.Client) *CmsApp {
	app := &CmsApp{
		db:  db,
		rdb: rdb,
	}
	return app
}
