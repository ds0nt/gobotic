package slack

import (
	"context"
	"strings"

	"github.com/ds0nt/gobotic/transports/types"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type SlackTransport struct {
	Token     string
	Channel   string
	onMessage types.MessageHandler
	onError   types.ErrorHandler
	logger    *logrus.Entry
	client    *slack.Client
	rtm       *slack.RTM
	ident     *slack.UserDetails
}

func NewSlackTransport(slackToken string, channel string) *SlackTransport {
	c := slack.New(slackToken)
	return &SlackTransport{
		client: c,
		rtm:    c.NewRTM(),
	}
}

func (t *SlackTransport) Connect(ctx context.Context) error {
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

func (t *SlackTransport) messageToTypesMessage(msg *slack.MessageEvent) (_msg types.MessageEvent) {
	_msg.Event = msg
	_msg.FullText = msg.Text

	if strings.HasPrefix(msg.Text, "<@"+t.ident.ID+"> ") {
		_msg.IsCommand = true
		_msg.ArgsText = strings.TrimPrefix(msg.Text, "<@"+t.ident.ID+"> ")
	}

	return
}

func (t *SlackTransport) SendMessage(ch string, msg string) {
	_, _, _, err := t.client.SendMessage(
		ch,
		slack.MsgOptionText(preWrap(msg), false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		t.logger.Errorln(ch, msg)
	}
}
func (t *SlackTransport) SendError(ch string, err error) {
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

func (t *SlackTransport) OnMessage(handler types.MessageHandler) {
	t.onMessage = handler
}

func (t *SlackTransport) OnError(handler types.ErrorHandler) {
	t.onError = handler
}

func (t *SlackTransport) Close() error {
	return t.rtm.Disconnect()
}

func (t *SlackTransport) Client() *slack.Client {
	return t.client
}

func (t *SlackTransport) RTM() *slack.RTM {
	return t.rtm
}

func (t *SlackTransport) Ident() *slack.UserDetails {
	return t.ident
}

func preWrap(s string) string {
	return "```\n" + s + "```"
}
