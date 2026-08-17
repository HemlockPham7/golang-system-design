package main

import (
	"context"

	workerHdl "github.com/HemlockPham7/golang-system-design/internal/handler/worker"
	"github.com/HemlockPham7/golang-system-design/internal/infrastructure"
	bookmarkRepo "github.com/HemlockPham7/golang-system-design/internal/repository/bookmark"
	"github.com/HemlockPham7/golang-system-design/internal/repository/queue"
	"github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	workerEng "github.com/HemlockPham7/golang-system-design/internal/worker"
	"github.com/HemlockPham7/golang-system-design/pkg/utils"
)

const queueName = "bookmark-queue"

func main() {
	// queue repository
	redisClient := infrastructure.CreateRedisClient("")
	queueRepository := queue.NewRedisQueue(redisClient, queueName)

	// db
	dbClient := infrastructure.CreateDB("")

	// handler
	bookmarkRepository := bookmarkRepo.NewRepository(dbClient)
	codeGen := utils.NewGenPass()
	bookmarkService := bookmark.NewService(bookmarkRepository, codeGen)
	handler := workerHdl.NewHandler(bookmarkService)

	// engine
	engine := workerEng.NewEngine(queueRepository, handler)
	engine.Start(context.Background())
}
