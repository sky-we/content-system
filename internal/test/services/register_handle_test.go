package services

import (
	"content-system/internal/services"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	// 启动go-mysql-server
	server, _, _ := GetMysqlFakeDBServer("cms_account", "account")
	go func() {
		if err := server.Start(); err != nil {
			fmt.Println("server start error", err)
			panic(err)
		}

	}()
	defer func() {
		if err := server.Close(); err != nil {
			fmt.Println("server Close error", err)
			panic(err)
		}
	}()

	// 启动MiniRedis
	s := miniredis.NewMiniRedis()

	if err := s.StartAddr("localhost:6380"); err != nil {
		t.Errorf("could not start miniredis with port 6380: %s", err)
		return
	}
	t.Cleanup(s.Close)

	// 路由绑定
	gin.SetMode(gin.TestMode)
	db := NewFakeMySqlDB(FakeDbCfg.MySQL)
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	app := services.NewCmsApp(db, rdb)
	r := gin.Default()
	r.POST("/register", app.Register)

	// 1 模拟正常请求，预期返回200，注册成功
	body := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	expectBody := `{"code":0,"msg":"ok","data":"haha register ok"}`
	assert.Equal(t, expectBody, w.Body.String())

	// 2 重复注册
	errorExistBody := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
	req2, err2 := http.NewRequest(http.MethodPost, "/register", strings.NewReader(errorExistBody))
	assert.NoError(t, err2)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)
	expectErrBody := "{\"Message\":\"账号已存在\"}"
	assert.Equal(t, expectErrBody, w2.Body.String())

	// 3 注册信息错误
	errorBody := `{"user_i":"haha","pass_word":"123456","nick_name":"sky-we"}`
	req, err3 := http.NewRequest(http.MethodPost, "/register", strings.NewReader(errorBody))
	assert.NoError(t, err3)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req)
	assert.Equal(t, 400, w3.Code)
	expectErrReqBody := `{"err":"Key: 'RegisterReq.UserId' Error:Field validation for 'UserId' failed on the 'required' tag"}`
	assert.Equal(t, expectErrReqBody, w3.Body.String())

}
