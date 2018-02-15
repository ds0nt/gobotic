package types

// MessageHandler is a message handler function.
type MessageHandler func(MessageEvent) error

// ErrorHandler is an error handler function.
type ErrorHandler func(error)

// MessageEvent is a message event.
type MessageEvent struct {
	FullText  string
	Channel   string
	IsCommand bool
	ArgsText  string
	Event     interface{}
}
