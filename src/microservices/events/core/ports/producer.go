package ports

import (
	"context"
)

type EventProducer interface {
	SendEvent(ctx context.Context, topic string, key []byte, value []byte) error
	Close()
}
