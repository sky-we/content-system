package services

import (
	"content-system/internal/services"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
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
	pro, _ := CreateTestDatabase(dbName, tableName, schema)
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
	r.POST("/register", app.Register)
	{
		// 1 模拟正常请求，预期返回200，注册成功
		body := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
		req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		expectBody := `{"code":0,"msg":"ok","data":"haha register ok"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

	{
		// 2 重复注册
		body := `{"user_id":"haha","pass_word":"123456","nick_name":"sky-we"}`
		req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectBody := `{"Message":"账号已存在"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

	{
		// 3 注册信息错误
		body := `{"user_i":"haha","pass_word":"123456","nick_name":"sky-we"}`
		req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectBody := `{"err":"Key: 'RegisterReq.UserId' Error:Field validation for 'UserId' failed on the 'required' tag"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

}
