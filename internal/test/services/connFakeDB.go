package services

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var FakeDbCfg *FakeDBConfig

type FakeMysqlConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	ChartSet    string
	ParseTime   string
	Loc         string
	MaxOpenConn int
	MaxIdleConn int
}

type FakeDBConfig struct {
	MySQL *FakeMysqlConfig
}

func LoadFakeDBConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../base/fakeConfig")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("viper ReadInConfig /base/fakeConfig panic")
		panic(err)

	}
	if err := viper.Unmarshal(&FakeDbCfg); err != nil {
	}

}

func NewFakeMySqlDB(cfg *FakeMysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.ChartSet,
		cfg.ParseTime,
		cfg.Loc,
	)
	mysqlDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		fmt.Println("connect mysql error:", err)
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		fmt.Println("get mysql instance error:", err)
		panic(err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	return mysqlDB
}
