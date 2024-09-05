package services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbe(t *testing.T) {
	r := gin.Default()
	cmsApp := NewCmsApp()
	r.GET("/api/cms/probe", cmsApp.Probe)
	url := "/api/cms/probe"
	Convey("Given a GET request to /api/cms/probe with session id", t, func() {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("session_id", "128sh123")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		Convey("When request with session id", func() {
			var rsp map[string]any
			So(rr.Code, ShouldEqual, http.StatusOK)
			err := json.Unmarshal(rr.Body.Bytes(), &rsp)
			Convey("response body json unmarshal should not throw error", func() {
				So(err, ShouldBeNil)
			})
			data := map[string]any{
				"Message": "Service Online!",
			}
			Convey("response body message should contain msg code data", func() {
				So(rsp["code"], ShouldEqual, 0)
				So(rsp["msg"], ShouldEqual, "ok")
				So(rsp["data"], ShouldEqual, data)

			})

		})

	})

}
