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

func (suite *ContentTestSuite) TestUpdateContentNoLogin() {
	reqBody := `{"id" : 6,
				"title":"Animal video2",
				"video_url":"www.baidu.com",
				"author":"sky-we",
				"description":"My first TikTok video",
				"thumbnail":"http://example.com/1.jpg",
				"category":"video",
				"duration":3600,
				"resolution":"720P",
				"file_size":10240,
				"format":"MP4",
				"quality":1,
				"approval_status":1
				}`

	// 未登录
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/update", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusUnauthorized, w.Code)
	expectBody := `{"Message":"用户未登录"}`
	suite.Equal(expectBody, w.Body.String())
}
func (suite *ContentTestSuite) TestUpdateContentArgsErr() {
	reqBody := `{
				"title":"cat video",
				"video_url":"www.sina1.com",
				"author":"sky-we",
				"description":"My second TikTok video",
				"thumbnail":"http://example.com/1.jpg",
				"category":"video",
				"duration":3600,
				"resolution":"720P",
				"file_size":10240,
				"format":"MP4",
				"quality":1,
				"approval_status":1
				}`
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/update", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"参数错误","err":"Key: 'ContentUpdateReq.ID' Error:Field validation for 'ID' failed on the 'required' tag"}`
	suite.Equal(expectBody, w.Body.String())
}
func (suite *ContentTestSuite) TestUpdateContentOk() {
	reqBody := `{"id" : 6,
				"title":"Animal video2",
				"video_url":"www.baidu.com",
				"author":"sky-we",
				"description":"My first TikTok video",
				"thumbnail":"http://example.com/1.jpg",
				"category":"video",
				"duration":3600,
				"resolution":"720P",
				"file_size":10240,
				"format":"MP4",
				"quality":1,
				"approval_status":1
				}`
	// 正常更新
	rowData := sql.Row{
		int32(6), "cat video", "My second TikTok video", "sky-we",
		"www.baidu.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)

	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/update", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"message":"success","data":{"ID":6}}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *ContentTestSuite) TestUpdateContentIDNotExist() {
	reqBody := `{"id" : 5,
				"title":"Animal video2",
				"video_url":"www.baidu.com",
				"author":"sky-we",
				"description":"My first TikTok video",
				"thumbnail":"http://example.com/1.jpg",
				"category":"video",
				"duration":3600,
				"resolution":"720P",
				"file_size":10240,
				"format":"MP4",
				"quality":1,
				"approval_status":1
				}`
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/update", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"[ID=5]内容不存在"}`
	suite.Equal(expectBody, w.Body.String())
}
