package services

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestContent 会运行所有绑定ContentTestSuite结构体的test
func TestContent(t *testing.T) {
	suite.Run(t, new(ContentTestSuite))
}
