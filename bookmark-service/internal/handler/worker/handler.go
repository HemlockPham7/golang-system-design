package worker

import (
	"context"
	"time"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, message []byte) error {
	println(string(message))
	time.Sleep(1 * time.Second)
	return nil
}
