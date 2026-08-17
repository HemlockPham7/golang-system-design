package worker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	"github.com/HemlockPham7/golang-system-design/internal/service/queue"
)

type Handler interface {
	Handle(ctx context.Context, message []byte) error
}

type handler struct {
	bookmarkService bookmark.Service
}

func NewHandler(bookmarkService bookmark.Service) Handler {
	return &handler{
		bookmarkService: bookmarkService,
	}
}

var ErrUnmarshalMessage = errors.New("failed to unmarshal message")

func (h *handler) Handle(ctx context.Context, message []byte) error {
	input := &queue.ImportMessage{}
	err := json.Unmarshal(message, input)
	if err != nil {
		return ErrUnmarshalMessage
	}

	err = h.bookmarkService.CreateBatchBookmarks(ctx, input.UID, input.Bookmarks)
	return nil
}
