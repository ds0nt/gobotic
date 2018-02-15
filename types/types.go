package types

type MessageHandler func(MessageEvent) error
type ErrorHandler func(Error)

type MessageEvent struct {
	FullText  string
	ArgsText  string
	InputText string
	Channel   string
	IsCommand bool
	Event     interface{}
}

type Error struct {
	Event *MessageEvent
	Err   error
}
