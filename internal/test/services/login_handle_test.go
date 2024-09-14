package services

import (
	"content-system/internal/services"
	"content-system/internal/utils"
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
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
	LoadFakeDBConfig()
	dbName := "cms_account"
	tableName := "account"
	schema := sql.Schema{
		{Name: "id", Type: types.Int32, Nullable: false, Source: tableName, PrimaryKey: true, Comment: "主键ID", AutoIncrement: true},
		{Name: "user_id", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "用户id"},
		{Name: "pass_word", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "密码"},
		{Name: "nick_name", Type: types.Text, Nullable: true, Source: tableName, PrimaryKey: false, Comment: "昵称"},
		{Name: "created_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
		{Name: "updated_at", Type: types.MustCreateDatetimeType(query.Type_DATETIME, 6), Nullable: false, Source: tableName},
	}
	pro, table := CreateTestDatabase(dbName, tableName, schema)
	server, _ := FakeMysqlServer(pro)
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
	{
		body := `{"user_id":"haha","pass_word":"123456"}`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "login ok")
	}

	{
		// 2 用户重复登录
		body := `{"user_id":"haha","pass_word":"123456"}`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		expectBody := `{"Message":"用户已登录"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

	{
		// 3 错误的请求体
		body := `{"user_":"haha","pass_word":"123456"}`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		expectBody := `{"error":"Key: 'LoginReq.UserId' Error:Field validation for 'UserId' failed on the 'required' tag"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

	{
		// 4 不存在的用户ID
		body := `{"user_id":"haha1","pass_word":"123456"}`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		expectBody := `{"Message":"请输入正确的用户ID"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

	{
		// 5 错误的用户密码
		rdb.Del(context.Background(), utils.GenSessionKey("haha")) // 退出登录
		body := `{"user_id":"haha","pass_word":"1234567"}`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		expectBody := `{"Message":"用户密码错误"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

}
