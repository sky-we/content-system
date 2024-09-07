package services

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	r := gin.Default()
	cmsApp := NewCmsApp()
	url := "/out/api/cms/register"
	r.POST(url, cmsApp.Register)
	for _, v := range r.Routes() {
		fmt.Printf("%s %s\n", v.Method, v.Path)
	}
	Convey("Given a GET request to /out/api/cms/register with correct body", t, func() {
		reqBody := `{"nick_name":"sky-we", "user_id":"007", "pass_word":"123456"}`
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		fmt.Println("req_body", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		Convey("When request ok", func() {
			var rsp map[string]any
			So(rr.Code, ShouldEqual, http.StatusOK)
			err := json.Unmarshal(rr.Body.Bytes(), &rsp)
			Convey("response body json unmarshal should not throw error", func() {
				So(err, ShouldBeNil)
			})
			Convey("response body message should contain msg code data", func() {
				So(rsp["code"], ShouldEqual, 0)
				So(rsp["msg"], ShouldEqual, "ok")
				So(rsp["data"], ShouldEqual, "register ok")

			})

		})

	})
	Convey("Given a GET request to /out/api/cms/register with error user_name", t, func() {
		reqBody := `{"user":"sky-we", "user_id":"007", "pass_word":"123456"}`
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		fmt.Println("req_body", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		Convey("When request ok", func() {
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
		})

	})

}
