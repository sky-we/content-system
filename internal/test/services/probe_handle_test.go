package services

import (
	"content-system/internal/api"
	"content-system/internal/services"
	"content-system/internal/utils"
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/dolthub/vitess/go/vt/proto/query"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProbe(t *testing.T) {
	LoadFakeDBConfig()
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

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	db := NewFakeMySqlDB(FakeDbCfg.MySQL)
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	cmsApp := services.NewCmsApp(db, rdb)
	sessionMiddleware := &api.SessionAuth{Rdb: rdb}
	root := r.Group(RootPath).Use(sessionMiddleware.Auth)
	root.GET("/cms/probe", cmsApp.Probe)

	{
		// 未登录
		req, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		expectBody := `{"Message":"用户未登录"}`
		assert.JSONEq(t, expectBody, w.Body.String())
	}
	{
		// 已登录
		request, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)
		assert.NoError(t, err)
		sessionId := uuid.New().String()
		request.Header.Set("session_id", sessionId)
		rdbErr := rdb.Set(context.Background(), utils.GenAuthKey(sessionId), time.Now().Unix(), time.Hour*8).Err()
		assert.NoError(t, rdbErr)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		assert.Equal(t, http.StatusOK, w.Code)
		expectBody := `{"code":0,"data":{"Message":"Service Online!"},"msg":"ok"}`
		assert.Equal(t, expectBody, w.Body.String())
	}

}
