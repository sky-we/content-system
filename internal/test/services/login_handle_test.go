package services

import (
	"context"
	"github.com/dolthub/go-mysql-server/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func (suite *AccountTestSuite) TestLoginOk() {
	passwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	rowData := sql.Row{int32(3), "haha", string(passwd), "sky-we", time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	body := `{"user_id":"haha","pass_word":"123456"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/login", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "login ok")
}

func (suite *AccountTestSuite) TestLoginRepeat() {
	_, err := suite.App.GenSessionId(context.Background(), "haha")
	suite.NoError(err, "GenSessionId should not throw error")
	body := `{"user_id":"haha","pass_word":"123456"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/login", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"用户已登录"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *AccountTestSuite) TestArgError() {
	body := `{"user_":"haha","pass_word":"123456"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/login", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"error":"Key: 'LoginReq.UserId' Error:Field validation for 'UserId' failed on the 'required' tag"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *AccountTestSuite) TestLoginNotExistUserID() {
	body := `{"user_id":"haha1","pass_word":"123456"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/login", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"请输入正确的用户ID"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *AccountTestSuite) TestErrorPasswd() {
	passwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	rowData := sql.Row{int32(3), "haha", string(passwd), "sky-we", time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	body := `{"user_id":"haha","pass_word":"1234567"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/login", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"用户密码错误"}`
	suite.Equal(expectBody, w.Body.String())
}
