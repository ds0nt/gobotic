package gobotic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"

	"github.com/ds0nt/gobotic/types"
	"github.com/pkg/errors"
)

type Command struct {
	Name    string
	Help    string
	Handler types.MessageHandler
}

type Interceptor func(msg types.MessageEvent) error

const CommandNameHelp = "help"

type CommandRouter struct {
	commandMap   map[string]*Command
	interceptors []Interceptor
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commandMap:   map[string]*Command{},
		interceptors: []Interceptor{},
	}
}

func (c *CommandRouter) Add(cmd *Command) {
	c.commandMap[cmd.Name] = cmd
}

func (c *CommandRouter) AddInterceptor(i Interceptor) {
	c.interceptors = append(c.interceptors, i)
}

func (c *CommandRouter) match(text string) (cmd *Command, input string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 1 {
		cmd, _ = c.commandMap[CommandNameHelp]
		return
	}

	cmd, ok := c.commandMap[parts[0]]
	if !ok {
		cmd, _ = c.commandMap[CommandNameHelp]
		return
	}

	if len(parts) == 2 {
		input = parts[1]
	}
	return
}

func (c *CommandRouter) Run(msg types.MessageEvent) error {
	for _, i := range c.interceptors {
		if err := i(msg); err != nil {
			return err
		}
	}
	cmd, input := c.match(msg.ArgsText)
	if cmd == nil {
		return errors.Errorf("command '%s' unknown", msg.ArgsText)
	}
	msg.InputText = input
	return cmd.Handler(msg)
}

func IsCommandNotFound(err error) bool {
	return strings.HasPrefix(err.Error(), "command '") && strings.HasSuffix(err.Error(), "' unknown")
}

var help = `%s bot usage:
	@%s commmand-name [optional command-input-string]...

Commands:
%s`

func (c *CommandRouter) Help(botName string) string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 4, 1, ' ', tabwriter.TabIndent)
	for _, c := range c.commandMap {
		fmt.Fprintf(w, "\t\t%s\t%s\n", c.Name, c.Help)
	}
	w.Flush()
	bytes, _ := ioutil.ReadAll(&buf)
	return fmt.Sprintf(help, botName, botName, string(bytes))
}
