package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// FakeMysqlConfig 模拟mysql go-mysql-server 配置
type FakeMysqlConfig struct {
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

// FakeRedisConfig 模拟redis Miniredis
type FakeRedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type FakeDBConfig struct {
	MySQL FakeMysqlConfig
	Redis FakeRedisConfig
}

func LoadFakeDBConfig() (*FakeDBConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/internal/test/config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading db config file, %s", err)

	}
	var fakeDbConfig FakeDBConfig
	if err := viper.Unmarshal(&fakeDbConfig); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)

	}
	return &fakeDbConfig, nil
}

var fakeDBConfig FakeDBConfig
