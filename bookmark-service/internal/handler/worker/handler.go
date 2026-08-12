package worker

import "context"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, message []byte) error {
	println(string(message))
	return nil
}
