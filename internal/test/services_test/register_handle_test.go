package services

import (
	"github.com/dolthub/go-mysql-server/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func (suite *AccountTestSuite) TestRegisterOk() {
	// 1 模拟正常请求，预期返回200，注册成功
	body := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/register", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusOK, w.Code)
	expectBody := `{"code":0,"msg":"ok","data":"haha register ok"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *AccountTestSuite) TestRegisterRepeat() {
	passwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	rowData := sql.Row{int32(3), "haha", string(passwd), "sky-we", time.Now(), time.Now()}
	InsertData(suite.DbName, suite.Provider, suite.Table, rowData)
	body := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/register", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"Message":"账号已存在"}`
	suite.Equal(expectBody, w.Body.String())
}

func (suite *AccountTestSuite) TestRegisterArgsError() {
	body := `{"user_i":"haha","pass_word":"123456","nick_name":"sky-we"}`
	request, err := http.NewRequest(http.MethodPost, "/out/api/cms/register", strings.NewReader(body))
	suite.NoError(err, "http.NewRequest should not throw error")
	w := httptest.NewRecorder()
	suite.GinEngine.ServeHTTP(w, request)
	suite.Equal(http.StatusBadRequest, w.Code)
	expectBody := `{"err":"Key: 'RegisterReq.UserId' Error:Field validation for 'UserId' failed on the 'required' tag"}`
	suite.Equal(expectBody, w.Body.String())
}
