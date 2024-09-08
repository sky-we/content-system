package database

import "github.com/stretchr/testify/mock"

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) name() {

}
