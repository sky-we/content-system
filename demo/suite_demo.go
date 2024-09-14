package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// MyTestSuite 定义了一个测试套件，其中包含一个变量
type MyTestSuite struct {
	suite.Suite
	// 定义一个变量，假设它是一个字符串
	testVariable string
}

// SetupSuite 在所有测试开始前执行，用于初始化那些所有测试共享的资源
func (suite *MyTestSuite) SetupSuite() {
	suite.testVariable = "initial value"
	suite.T().Log("SetupSuite is invoked")
}

// SetupTest 在每个测试方法前执行，用于初始化每个测试方法的资源
func (suite *MyTestSuite) SetupTest() {
	suite.T().Log("SetupTest is invoked")
}

// TearDownTest 在每个测试方法后执行，用于清理每个测试方法的资源
func (suite *MyTestSuite) TearDownTest() {
	suite.T().Log("TearDownTest is invoked")
}

// TearDownSuite 在所有测试完成后执行，用于清理测试套件的资源
func (suite *MyTestSuite) TearDownSuite() {
	suite.T().Log("TearDownSuite is invoked")
}

// TestExample1 是一个测试方法，使用定义的变量
func (suite *MyTestSuite) TestExample1() {
	suite.Equal("initial value", suite.testVariable, "testVariable should be 'initial value'")
	suite.T().Log("TestExample1 is running, testVariable =", suite.testVariable)
}

// TestExample2 是另一个测试方法，也使用定义的变量
func (suite *MyTestSuite) TestExample2() {
	suite.testVariable = "modified value"
	suite.Equal("modified value", suite.testVariable, "testVariable should be 'modified value'")
	suite.T().Log("TestExample2 is running, testVariable =", suite.testVariable)
}

// TestMyTestSuite 用于注册测试套件
func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
