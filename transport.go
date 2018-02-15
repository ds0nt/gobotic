package gobotic

import (
	"context"

	"github.com/ds0nt/gobotic/types"
)

type Transport interface {
	Connect(ctx context.Context) error
	OnMessage(types.MessageHandler)
	OnError(types.ErrorHandler)
	BotID() string
	BotName() string
	Send(string, string)
	Close() error
}
