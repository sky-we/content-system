package services

import (
	"content-system/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
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
	connDB(app)
	connRdb(app)
	return app
}

func NewMysqlDB(cfg config.MysqlConfig) *gorm.DB {
	//dsn := "root:1qaz!QAZ@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	return mysqlDB
}

func NewRdb(option redis.Options) *redis.Client {

	rdb := redis.NewClient(&option)
	return rdb
}
