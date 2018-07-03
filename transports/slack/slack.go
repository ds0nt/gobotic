package slack

import (
	"context"
	"strings"

	"github.com/ds0nt/gobotic/types"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

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

func NewTransport(slackToken string, channel string, logger *logrus.Entry) *Transport {
	c := slack.New(slackToken)
	if logger == nil {
		logger = logrus.NewEntry(logrus.New())
	}
	return &Transport{
		client: c,
		rtm:    c.NewRTM(),
		logger: logger,
	}
}

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
								t.onError(types.Error{
									Err:   err,
									Event: &e,
								})
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
	_msg.Channel = msg.Channel
	_msg.User = msg.User
	_msg.Transport = t

	if strings.HasPrefix(msg.Text, "<@"+t.ident.ID+"> ") {
		_msg.IsCommand = true
		_msg.ArgsText = strings.TrimPrefix(msg.Text, "<@"+t.ident.ID+"> ")
	}

	return
}

func (t *Transport) Send(channel string, text string) {
	_, _, _, err := t.client.SendMessage(
		channel,
		slack.MsgOptionText(PreWrap(text), false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		t.logger.Errorln(channel, text, err)
	}
}

func PreWrap(s string) string {
	return "```\n" + s + "```"
}

func (t *Transport) OnMessage(handler types.MessageHandler) {
	t.onMessage = handler
}

func (t *Transport) OnError(handler types.ErrorHandler) {
	t.onError = handler
}

func (t *Transport) Close() error {
	return t.rtm.Disconnect()
}

func (t *Transport) Client() *slack.Client {
	return t.client
}

func (t *Transport) RTM() *slack.RTM {
	return t.rtm
}

func (t *Transport) Ident() *slack.UserDetails {
	return t.ident
}

func (t *Transport) BotID() string {
	if t.ident == nil {
		return ""
	}
	return t.ident.ID
}

func (t *Transport) BotName() string {
	if t.ident == nil {
		return ""
	}
	return t.ident.Name
}
