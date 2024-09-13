package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type MysqlConfig struct {
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

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type RedisWinConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type DataBaseConfig struct {
	MySQL    *MysqlConfig
	Redis    *RedisConfig
	RedisWin *RedisWinConfig
}

var (
	once     sync.Once
	DBConfig *DataBaseConfig
)

func LoadDBConfig() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("d://Code/content-system/internal/config")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("error reading db config file, %s", err)
			panic(err)

		}
		if err := viper.Unmarshal(&DBConfig); err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
			panic(err)
		}
	})
}

func NewMySqlDB(cfg *MysqlConfig) *gorm.DB {
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
	fmt.Println("mysql connect dsn:", dsn)
	mysqlDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	return mysqlDB
}

func NewRdb(cfg *RedisConfig) *redis.Client {
	option := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	fmt.Println("redis connect option:", option)
	rdb := redis.NewClient(&option)
	return rdb
}
