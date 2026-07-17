package main

import (
	"context"
	"time"

	"github.com/HemlockPham7/golang-system-design/pkg/redis"
)

func main() {
	rclient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	rclient.Set(context.Background(), "12345", "google.com", time.Hour)

	rclient2, err := redis.NewClient("CACHE")
	if err != nil {
		panic(err)
	}

	rclient2.Set(context.Background(), "31726", "500 days", time.Hour)
}
