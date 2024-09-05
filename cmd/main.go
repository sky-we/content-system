package main

import (
	"content-system/internal/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.CmsRouters(r)
	if err := r.Run(); err != nil {
		fmt.Printf("run err %v", err)
		return
	}

}
