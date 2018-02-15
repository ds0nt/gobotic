package gobotic

import (
	"context"

	"github.com/ds0nt/gobotic/transports/types"
)

// Command is a bot command.
type Command struct {
	Name    string
	Help    string
	Handler func(msg types.MessageEvent, input string) error
}

// Interceptor is a function to intercept a message.
type Interceptor func(msg types.MessageEvent) error

// Transport defines a standard transport layer.
type Transport interface {
	Connect(ctx context.Context) error
	OnMessage(types.MessageHandler)
	OnError(types.ErrorHandler)
	SendMessage(string, string)
	SendError(string, error)
	Close() error
}
