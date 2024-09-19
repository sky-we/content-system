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

func (suite *ContentTestSuite) TestFindContentNotLogin() {

	reqBody := `{"author" :"sky-we","page":2,"pageSize":1}`
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/find", strings.NewReader(reqBody))
	suite.NoError(err)
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusUnauthorized, w.Code)
	expectBody := `{"Message":"用户未登录"}`
	suite.Equal(expectBody, w.Body.String())

}

func (suite *ContentTestSuite) TestFindContentByAuthor() {
	rowData := sql.Row{
		int32(1), "cat video", "My second TikTok video", "sky-we",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData2 := sql.Row{
		int32(2), "cat video", "My second TikTok video", "sky-we",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData3 := sql.Row{
		int32(3), "cat video", "My second TikTok video", "sky-we2",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData4 := sql.Row{
		int32(4), "cat video", "My second TikTok video", "sky-we2",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData2)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData3)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData4)
	reqBody := `{"author" :"sky-we","page":2,"pageSize":1}`
	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/find", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"message":"success","data":{"content":[{"ID":2,"Title":"cat video","Description":"My second TikTok video","Author":"sky-we","VideoURL":"www.sina1.com","Thumbnail":"http://example.com/1.jpg","Category":"video","Duration":3600,"Resolution":"720P","FileSize":10240,"Format":"MP4","Quality":1,"ApprovalStatus":1}],"total":2}}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *ContentTestSuite) TestFindContentByTitle() {
	rowData := sql.Row{
		int32(1), "cat video", "My second TikTok video", "sky-we",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData2 := sql.Row{
		int32(2), "cat video", "My second TikTok video", "sky-we",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData3 := sql.Row{
		int32(3), "dog video", "My second TikTok video", "sky-we2",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	rowData4 := sql.Row{
		int32(4), "cat video", "My second TikTok video", "sky-we2",
		"www.sina1.com", "http://example.com/1.jpg",
		"video", int64(3600), "720P", int64(10240), "MP4", int8(1), int8(1), time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData2)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData3)
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData4)

	reqBody := `{"title" :"cat video","page":2,"pageSize":1}`

	request, err := http.NewRequest(http.MethodPost, "/api/cms/content/find", strings.NewReader(reqBody))
	suite.NoError(err, "http.NewRequest should not throw error")
	sessionId := uuid.New().String()
	request.Header.Set("session_id", sessionId)
	rdbErr := suite.Rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
	suite.NoError(rdbErr, "rdb.Set should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"message":"success","data":{"content":[{"ID":2,"Title":"cat video","Description":"My second TikTok video","Author":"sky-we","VideoURL":"www.sina1.com","Thumbnail":"http://example.com/1.jpg","Category":"video","Duration":3600,"Resolution":"720P","FileSize":10240,"Format":"MP4","Quality":1,"ApprovalStatus":1}],"total":3}}`
	suite.Equal(expectBody, w.Body.String())
}
