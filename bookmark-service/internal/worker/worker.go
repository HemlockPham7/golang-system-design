package worker

import (
	"context"

	"github.com/rs/zerolog/log"
)

type pool struct { // quan li so luong worker trong mot pool
	handler      Handler
	numberWorker int
	messages     chan []byte
}

func newPool(ctx context.Context, handler Handler, numberWorker int) *pool {
	messageChan := make(chan []byte, numberWorker)

	initPool := &pool{
		handler:      handler,
		numberWorker: numberWorker,
		messages:     messageChan,
	}

	initPool.init(ctx)
	return initPool
}

func (p *pool) init(ctx context.Context) {
	for i := 0; i < p.numberWorker; i++ {
		w := &worker{
			id:       i + 1,
			handler:  p.handler,
			messages: p.messages,
		}
		log.Info().Msgf("Starting worker %d", w.id)
		go w.run(ctx)
	}
}

func (p *pool) Consume(message []byte) {
	p.messages <- message
}

type worker struct {
	id       int
	handler  Handler
	messages <-chan []byte
}

func (w *worker) run(ctx context.Context) {
	for {
		msg, ok := <-w.messages
		if !ok {
			log.Info().Msgf("Worker %d is closing", w.id)
			return
		}
		log.Debug().Msgf("Worker %d is processing message: %s", w.id, string(msg))
		err := w.handler.Handle(ctx, msg)
		if err != nil {
			log.Error().Err(err).Msgf("Worker %d failed to process message", w.id)
		} else {
			log.Info().Msgf("Worker %d processed successfully message: %s", w.id, string(msg))
		}
	}
}
