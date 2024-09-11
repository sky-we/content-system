package services

import (
	"content-system/internal/services"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	// 启动go-mysql-server
	server, pro, table := GetMysqlFakeDBServer("cms_account", "account")
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
	r.POST("/login", app.Login)

	// 数据准备
	passwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	rowData := sql.Row{int32(3), "haha", string(passwd), "sky-we", time.Now(), time.Now()}
	InsertData("cms_account", pro, table, rowData)
	// 1 模拟正常请求，预期返回200，登录成功
	body := `{"user_id":"haha","pass_word":"123456"}`
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

}
