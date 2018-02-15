package slack

import (
	"context"
	"strings"

	"github.com/ds0nt/gobotic/transports/types"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

// Transport provides transport access to slack.
type Transport struct {
	Token     string
	Channel   string
	onMessage types.MessageHandler
	onError   types.ErrorHandler
	logger    *logrus.Entry
	client    *slack.Client
	rtm       *slack.RTM
	ident     *slack.UserDetails
}

// NewTransport returns a new instance of Transport.
func NewTransport(token string, channel string) *Transport {
	c := slack.New(token)
	return &Transport{
		client: c,
		rtm:    c.NewRTM(),
	}
}

// Connect initializes our transport connection.
func (t *Transport) Connect(ctx context.Context) error {
	go t.rtm.ManageConnection()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-t.rtm.IncomingEvents:
				func() {
					defer func() {
						if err := recover(); err != nil {
							t.logger.Errorf("%v", err)
						}
					}()

					switch event := msg.Data.(type) {
					case *slack.ConnectedEvent:
						t.ident = event.Info.User
						t.logger.Infof("Connected %v", event)

					case *slack.MessageEvent:
						go func(e types.MessageEvent) {
							err := t.onMessage(e)
							if err != nil {
								t.onError(err)
							}
						}(t.messageToTypesMessage(event))
					case *slack.MessageTooLongEvent:
						t.logger.Errorf("%v", event)
					case *slack.ReconnectUrlEvent:
						t.logger.Infof("%v", event)
					case *slack.DisconnectedEvent:
						t.logger.Infof("%v", event)
					default:
						t.logger.Infof("%v", event)
					}
				}()
			}
		}
	}()

	return nil
}

func (t *Transport) messageToTypesMessage(msg *slack.MessageEvent) (_msg types.MessageEvent) {
	_msg.Event = msg
	_msg.FullText = msg.Text

	if strings.HasPrefix(msg.Text, "<@"+t.ident.ID+"> ") {
		_msg.IsCommand = true
		_msg.ArgsText = strings.TrimPrefix(msg.Text, "<@"+t.ident.ID+"> ")
	}

	return
}

// SendMessage sends a message over the transport.
func (t *Transport) SendMessage(ch string, msg string) {
	_, _, _, err := t.client.SendMessage(
		ch,
		slack.MsgOptionText(preWrap(msg), false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		t.logger.Errorln(ch, msg)
	}
}

// SendError sends an error over the transport.
func (t *Transport) SendError(ch string, err error) {
	_, _, _, err = t.client.SendMessage(
		ch,
		slack.MsgOptionText(preWrap(err.Error()), false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		t.logger.Errorln(ch, err)
		return
	}
}

// OnMessage defines how to react when receiving a message over the transport.
func (t *Transport) OnMessage(handler types.MessageHandler) {
	t.onMessage = handler
}

// OnError defines how to react when receiving an error over the transport.
func (t *Transport) OnError(handler types.ErrorHandler) {
	t.onError = handler
}

// Close closes the transport.
func (t *Transport) Close() error {
	return t.rtm.Disconnect()
}

// Client returns any underlying transport client.
func (t *Transport) Client() *slack.Client {
	return t.client
}

// RTM returns the underlying RTM.
func (t *Transport) RTM() *slack.RTM {
	return t.rtm
}

// Ident returns the transports logged in user's identification.
func (t *Transport) Ident() *slack.UserDetails {
	return t.ident
}

func preWrap(s string) string {
	return "```\n" + s + "```"
}
