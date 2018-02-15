package gobotic

import (
	"context"

	"github.com/ds0nt/gobotic/transports/types"
)

// Bot is a bot, is a bot.
type Bot struct {
	transport Transport
	router    *CommandRouter
}

// NewBot returns a bot using a specified transport and router.
func NewBot(t Transport, r *CommandRouter) *Bot {
	return &Bot{
		transport: t,
		router:    r,
	}
}

// Run runs the bot.
func (c *Bot) Run(ctx context.Context) error {
	err := c.transport.Connect(ctx)
	if err != nil {
		return err
	}
	c.transport.OnMessage(c.OnMessage)
	//c.transport.OnError(c.OnError)

	return nil
}

// OnMessage defines the bots response to a message event.
func (c *Bot) OnMessage(msg types.MessageEvent) error {
	if msg.IsCommand {
		err := c.router.Run(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// OnError defines the bots response to an error event.
// func (c *Bot) OnError(err error) {
// 	c.transport.SendError(err)
// }
