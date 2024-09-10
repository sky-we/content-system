package services

import (
	"content-system/internal/services"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	go func() {
		startMysqlFakeDB()
	}()
	go func() {
		s := miniredis.NewMiniRedis()
		// redis测试端口6380
		if err := s.StartAddr("localhost:6380"); err != nil {
			t.Errorf("could not start miniredis: %s", err)
			return
		}
		t.Cleanup(s.Close)
	}()
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	db := NewFakeMySqlDB(FakeDbCfg.MySQL)

	rdb := NewFakeRdb(FakeDbCfg.Redis)

	app := services.NewCmsApp(db, rdb)

	r.POST("/register", app.Register)

	body := `{"user_id":"haha","pass_word":"123456","nick_name":"lwlw"}`

	req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(body))

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expectBody := `{"code":0,"msg":"ok","data":"haha register ok"}`

	assert.Equal(t, expectBody, w.Body.String())

}
