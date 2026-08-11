package queue

import "context"

//go:generate mockery --name Repository --filename common.go --outpkg mock_queue
type Repository interface {
	PushMessage(ctx context.Context, message []byte) error
}
