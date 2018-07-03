package types

import (
	"context"
)

type MessageHandler func(MessageEvent) error
type ErrorHandler func(Error)

type MessageEvent struct {
	FullText  string
	ArgsText  string
	InputText string
	Channel   string
	IsCommand bool
	Event     interface{}
	User      string
	Transport Transport
}

type Error struct {
	Event *MessageEvent
	Err   error
}

type Transport interface {
	Connect(ctx context.Context) error
	OnMessage(MessageHandler)
	OnError(ErrorHandler)
	BotID() string
	BotName() string
	Send(string, string)
	Close() error
}
