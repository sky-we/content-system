package services

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CmsApp struct {
	db *gorm.DB
}

func NewCmsApp() *CmsApp {
	app := &CmsApp{}
	connDB(app)
	return app
}

func connDB(app *CmsApp) {
	dsn := "root:1qaz!QAZ@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
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
	app.db = mysqlDB
}
