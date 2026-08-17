package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type pool struct { // quan li so luong worker trong mot pool
	handler      Handler
	numberWorker int
	messages     chan []byte
	errChan      chan *worker
	wg           *sync.WaitGroup
}

func newPool(ctx context.Context, handler Handler, numberWorker int) *pool {
	messageChan := make(chan []byte, numberWorker)
	errorChan := make(chan *worker, numberWorker)

	initPool := &pool{
		handler:      handler,
		numberWorker: numberWorker,
		messages:     messageChan,
		errChan:      errorChan,
		wg:           &sync.WaitGroup{},
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
			errChan:  p.errChan,
			wg:       p.wg,
		}
		log.Info().Msgf("Starting worker %d", w.id)
		p.wg.Add(1)
		go w.run(ctx)
	}

	go func() {
		for w := range p.errChan {
			log.Error().Msgf("Worker %d exited with error %v", w.id, w.err)

			time.Sleep(1 * time.Second)
			log.Info().Msgf("Restarting worker %d ...", w.id)
			w.err = nil
			go w.run(ctx)
		}
	}()
}

func (p *pool) Consume(message []byte) {
	p.messages <- message
}

func (p *pool) Close() {
	close(p.messages)
	close(p.errChan)
	p.wg.Wait()
	log.Info().Msg("worker pool closed")
}

type worker struct {
	id       int
	handler  Handler
	messages <-chan []byte
	err      error
	errChan  chan *worker
	wg       *sync.WaitGroup
}

func (w *worker) run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				w.err = err
			} else {
				w.err = fmt.Errorf("panic happened with %v", r)
			}
			w.errChan <- w
		} else {
			w.wg.Done()
		}
	}()

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
