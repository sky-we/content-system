package main

import (
	"content-system/internal/api"
	"content-system/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func init() {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		panic(err)
	}

}

func main() {
	r := gin.Default()
	api.CmsRouters(r)
	if err := r.Run(); err != nil {
		fmt.Printf("run err %v", err)
		return
	}

}
