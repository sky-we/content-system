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
	"testing"
)

func TestProbe(t *testing.T) {
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

	gin.SetMode(gin.TestMode)

	r := gin.Default()

	db := NewFakeMySqlDB(FakeDbCfg.MySQL)

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	cmsApp := services.NewCmsApp(db, rdb)

	r.GET("/api/cms/probe", cmsApp.Probe)

	req, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectBody := `{"code":0,"msg":"ok","data":{"Message":"Service Online!"}}`

	assert.JSONEq(t, expectBody, w.Body.String())
}
