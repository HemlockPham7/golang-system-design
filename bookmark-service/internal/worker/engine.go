package worker

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/repository/queue"
	"github.com/rs/zerolog/log"
)

type Engine interface {
	Start(ctx context.Context)
}

type Handler interface {
	Handle(ctx context.Context, message []byte) error
}

type engine struct {
	queue   queue.Repository
	handler Handler
	run     bool
	sigChan chan os.Signal
}

func NewEngine(queue queue.Repository, handler Handler) Engine {
	return &engine{
		queue:   queue,
		handler: handler,
		run:     false,
		sigChan: make(chan os.Signal, 1),
	}
}

const (
	intervalDelay  = 1 * time.Second
	numberOfWorker = 4
)

func (e *engine) Start(ctx context.Context) {
	log.Info().Msg("Starting worker engine")

	workerPool := newPool(ctx, e.handler, numberOfWorker)
	signal.Notify(e.sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	e.run = true
	for e.run {
		select {
		case sig := <-e.sigChan:
			log.Info().Msgf("Received signal: %s", sig.String())
			e.run = false
		default:
			// pop message
			msg, err := e.queue.PopMessage(ctx)
			if err != nil {
				if errors.Is(err, queue.NoMessageError) {
					time.Sleep(intervalDelay)
					continue
				}

				log.Error().Err(err).Msg("Failed to pop message")
				time.Sleep(intervalDelay)
				continue
			}

			// handle message
			workerPool.Consume(msg)
		}
	}
	workerPool.Close()
}
