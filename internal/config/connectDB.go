package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type MysqlConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	ChartSet    string
	ParseTime   bool
	Loc         string
	MaxOpenConn int
	MaxIdleConn int
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

// RedisWinConfig TODO： window开发暂时用该配置，部署到linux环境后删除
type RedisWinConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type DBConfig struct {
	MySQL    MysqlConfig
	Redis    RedisConfig
	RedisWin RedisWinConfig
}

var (
	once     sync.Once
	dbConfig *DBConfig
)

func LoadDBConfig() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/internal/config")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("error reading db config file, %s", err)
			panic(err)

		}
		if err := viper.Unmarshal(&dbConfig); err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
			panic(err)
		}
	})
}
