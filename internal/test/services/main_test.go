package services

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestAll 单元测试入口 go test -run TestAll
func TestAll(t *testing.T) {
	t.Log("test all")
	// logrus是go-mysql-server依赖的日志包，设置为Error级别
	log.SetLevel(log.ErrorLevel)
	suite.Run(t, new(ContentTestSuite))
	suite.Run(t, new(AccountTestSuite))

}
