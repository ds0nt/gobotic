package gobotic

import (
	"context"

	"github.com/ds0nt/gobotic/transports/types"
)

type Command struct {
	Name    string
	Help    string
	Handler func(msg types.MessageEvent, input string) error
}

type Interceptor func(msg types.MessageEvent) error

type Transport interface {
	Connect(ctx context.Context) error
	OnMessage(types.MessageHandler)
	OnError(types.ErrorHandler)
	SendMessage(string)
	SendError(error)
	Close() error
}
