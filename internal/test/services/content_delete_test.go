package services

import (
	"content-system/internal/utils"
	"context"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func (suite *ContentTestSuite) TestUserNotLogin() {
	suite.Root.POST("/cms/content/delete", suite.App.ContentDelete)

	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/delete", strings.NewReader(`{"id" :6}`))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusUnauthorized, w.Code)
	expectBody := `{"Message":"用户未登录"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *ContentTestSuite) TestDeleteNotExistID() {
	suite.Root.POST("/cms/content/delete", suite.App.ContentDelete)

	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/delete", strings.NewReader(`{"id" :7}`))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"[ID=7]内容不存在]"}`
	suite.Equal(expectBody, w.Body.String())

}

func (suite *ContentTestSuite) TestArgsError() {
	suite.Root.POST("/cms/content/delete", suite.App.ContentDelete)

	reqBody := `{"id111":5}`
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/delete", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.T().Log("argserror routes", suite.GinEngine.Routes())
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"参数错误","error":"Key: 'ContentDeleteReq.ID' Error:Field validation for 'ID' failed on the 'required' tag"}`
	suite.Equal(expectBody, w.Body.String())

}

func (suite *ContentTestSuite) TestDeleteOk() {
	suite.Root.POST("/cms/content/delete", suite.App.ContentDelete)

	rowData := sql.Row{
		int32(6), "cat video", "My second TikTok video", "sky-we",
		"www.baidu.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/delete", strings.NewReader(`{"id":6}`))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.T().Log("TestDeleteOk routes", suite.GinEngine.Routes())

	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"msg":"success","data":"ID 6 delete"}`
	suite.Equal(expectBody, w.Body.String())
}
