# gobotic

Gopher bots


# Example Use


```go
package main

import (
	"context"
	"flag"

	"github.com/abadojack/whatlanggo"
	"github.com/ds0nt/gobotic"
	"github.com/ds0nt/gobotic/transports/slack"
	"github.com/ds0nt/gobotic/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	slackToken = flag.String("slack_token", "", "token for slack")
	log        = logrus.NewEntry(logrus.New())
	t          *slack.Transport
)

func main() {
	flag.Parse()

	t = slack.NewTransport(*slackToken, "", log)
	r := gobotic.NewCommandRouter()

	r.AddInterceptor(frenchInterceptor)
	r.Add(&gobotic.Command{
		Name:    "say",
		Help:    "repeats text",
		Handler: sayHandler,
	})
	r.Add(&gobotic.Command{
		Name:    "detect-language",
		Help:    "detect's language",
		Handler: detectLanguageHandler,
	})

	bot := gobotic.NewBot(t, r)
	ctx := context.Background()
	err := bot.Run(ctx)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()

}

func sayHandler(msg types.MessageEvent) error {
	t.Send(msg.Channel, msg.InputText)
	return nil
}
func detectLanguageHandler(msg types.MessageEvent) error {
	l := whatlanggo.DetectLang(msg.FullText)
	t.Send(msg.Channel, whatlanggo.Langs[l])
	return nil
}

func frenchInterceptor(msg types.MessageEvent) error {
	l := whatlanggo.DetectLang(msg.FullText)
	if l == whatlanggo.Fra {
		return errors.Errorf("Sorry, I don't execute les commands, mon ami. .")
	}
	return nil
}


```