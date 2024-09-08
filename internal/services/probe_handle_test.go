package services

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	cmsApp := NewCmsApp()

	r.GET("/api/cms/probe", cmsApp.Probe)

	req, err := http.NewRequest(http.MethodGet, "/api/cms/probe", nil)

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectBody := `{"code":0,"msg":"ok","data":{"Message":"Service Online!"}}`

	assert.JSONEq(t, expectBody, w.Body.String())
}
