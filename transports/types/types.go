package types

type MessageHandler func(MessageEvent) error
type ErrorHandler func(error)

type MessageEvent struct {
	FullText  string
	Channel   string
	IsCommand bool
	ArgsText  string
	Event     interface{}
}
