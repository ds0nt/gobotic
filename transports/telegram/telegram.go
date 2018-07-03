package telegram

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ds0nt/gobotic/types"
	"github.com/kr/pretty"
	"github.com/labstack/echo"
)

// Transport is the telegram transport
type Transport struct {
	webhookURL string
	httpAddr   string

	onMessage types.MessageHandler
	onError   types.ErrorHandler

	botToken   string
	telegram   *telegramClient
	httpServer *http.Server
}

func NewTransport(httpAddr, webhookURL, botToken string) *Transport {
	return &Transport{
		telegram:   newTelegramClient(botToken),
		webhookURL: webhookURL,
		httpAddr:   httpAddr,
		botToken:   botToken,
	}
}

func (t *Transport) Connect(ctx context.Context) error {
	err := t.telegram.SetWebhook(t.webhookURL)
	if err != nil {
		return err
	}
	e := echo.New()
	e.POST("/*", func(c echo.Context) (err error) {
		r := Update{}
		if err = c.Bind(&r); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if r.Message != nil {
			t.onMessage(t.messageToTypesMessage(r.Message))
		}

		return c.JSON(http.StatusOK, nil)
	})

	t.httpServer = &http.Server{
		Addr:    t.httpAddr,
		Handler: e,
	}
	go t.httpServer.ListenAndServe()
	go func() {
		<-ctx.Done()
		t.httpServer.Close()
	}()

	return nil
}

func (t *Transport) messageToTypesMessage(msg *Message) (_msg types.MessageEvent) {
	_msg.Event = msg
	_msg.FullText = msg.Text
	_msg.Channel = strconv.Itoa(int(msg.Chat.ID))
	_msg.User = msg.From.UserName
	_msg.Transport = t
	pretty.Println(msg)
	prefix := "/phoenix "

	if strings.HasPrefix(strings.ToLower(msg.Text), prefix) {
		_msg.IsCommand = true
		_msg.ArgsText = strings.TrimPrefix(strings.ToLower(msg.Text), prefix)
	}

	return
}

func (t *Transport) OnMessage(handler types.MessageHandler) {
	t.onMessage = handler
}

func (t *Transport) OnError(handler types.ErrorHandler) {
	t.onError = handler
}

func (t *Transport) BotID() string {
	return ""
}
func (t *Transport) BotName() string {
	return ""
}

func PreWrap(s string) string {
	return "```\n" + s + "```"
}
func (t *Transport) Send(channel, text string) {
	chatID, _ := strconv.Atoi(channel)
	t.telegram.SendMessage(&SendMessageRequest{
		ChatID: int64(chatID),
		Text:   PreWrap(text),
	})
	return
}
func (t *Transport) Close() error {
	t.httpServer.Close()
	return nil
}
