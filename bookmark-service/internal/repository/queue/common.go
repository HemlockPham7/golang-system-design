package queue

import "context"

type Repository interface {
	PushMessage(ctx context.Context, message []byte) error
}
