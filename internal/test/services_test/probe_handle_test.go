package services

import (
	"content-system/internal/utils"
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"time"
)

// Probe较简单 绑定在ContentTestSuite即可
func (suite *ContentTestSuite) TestProbeNoLogin() {
	request, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)
	suite.NoError(err, "GenSessionId should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusUnauthorized, w.Code)
	expectBody := `{"Message":"用户未登录"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *ContentTestSuite) TestProbeLoginAlready() {
	request, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)
	suite.NoError(err, "GenSessionId should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"data":{"Message":"Service Online!"},"msg":"ok"}`
	suite.Equal(expectBody, w.Body.String())
}
