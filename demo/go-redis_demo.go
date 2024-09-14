package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	connDB()
}
func connDB() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := rdb.Set(ctx, "lw", "Rich", 5*time.Second).Err(); err != nil {
		panic(err)
	}
	val, err := rdb.Get(ctx, "lw").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
