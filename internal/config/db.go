package config

import (
	"content-system/internal/process"
	"fmt"
	"github.com/redis/go-redis/v9"
	goflow "github.com/s8sg/goflow/v1"
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

type FlowServiceConfig struct {
	RedisURL          string
	Port              int
	WorkerConcurrency int
}

type DataBaseConfig struct {
	MySQL       *MysqlConfig
	Redis       *RedisConfig
	RedisWin    *RedisWinConfig
	FlowService *FlowServiceConfig
}

var (
	once     sync.Once
	DBConfig *DataBaseConfig
)

func LoadDBConfig() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("internal/config")

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

func NewFlowService(cfg *FlowServiceConfig) *goflow.FlowService {
	fs := goflow.FlowService{
		Port:              cfg.Port,
		RedisURL:          cfg.RedisURL,
		WorkerConcurrency: cfg.WorkerConcurrency,
	}
	contentFlow := process.ContentFlow{}
	err := fs.Register("content-flow", contentFlow.ContentFlowHandle)
	if err != nil {
		panic(err)
	}

	return &fs
}
