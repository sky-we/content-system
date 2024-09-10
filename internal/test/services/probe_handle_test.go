package services

import (
	"content-system/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbe(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.Default()

	db := NewFakeMySqlDB(FakeDbCfg.MySQL)

	rdb := NewFakeRdb(FakeDbCfg.Redis)

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
