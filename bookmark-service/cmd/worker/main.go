package main

import (
	"context"

	workerHdl "github.com/HemlockPham7/golang-system-design/internal/handler/worker"
	"github.com/HemlockPham7/golang-system-design/internal/infrastructure"
	"github.com/HemlockPham7/golang-system-design/internal/repository/queue"
	workerEng "github.com/HemlockPham7/golang-system-design/internal/worker"
)

const queueName = "bookmark-queue"

func main() {
	// queue repository
	redisClient := infrastructure.CreateRedisClient("")
	queueRepository := queue.NewRedisQueue(redisClient, queueName)

	// handler
	handler := workerHdl.NewHandler()

	// engine
	engine := workerEng.NewEngine(queueRepository, handler)
	engine.Start(context.Background())
}
