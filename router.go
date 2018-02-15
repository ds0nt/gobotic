package gobotic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"

	"github.com/ds0nt/gobotic/transports/types"
)

// CommandNameHelp probably doesn't need to be externalized...
const CommandNameHelp = "help"

// CommandRouter matches input strings against known commands and routes
// them to the appropriate functions.
type CommandRouter struct {
	commandMap   map[string]*Command
	interceptors []Interceptor
}

// NewCommandRouter returns a new CommandRouter
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commandMap:   map[string]*Command{},
		interceptors: []Interceptor{},
	}
}

// Add a command to the router.
func (c *CommandRouter) Add(cmd *Command) {
	c.commandMap[cmd.Name] = cmd
}

// AddInterceptor adds an interceptor to the router.
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

// Run should be called when a new message event needs to be processed
// by the router.
func (c *CommandRouter) Run(msg types.MessageEvent) error {
	for _, i := range c.interceptors {
		if err := i(msg); err != nil {
			return err
		}
	}
	cmd, input := c.match(msg.ArgsText)
	return cmd.Handler(msg, input)
}

var help = `%s bot usage:
	@%s commmand-name [optional command-input-string]...

Commands:
%s`

// Help returns the help string for the given command.
func (c *CommandRouter) Help(id string) string {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 4, 1, ' ', tabwriter.TabIndent)
	for _, c := range c.commandMap {
		fmt.Fprintf(w, "\t\t%s\t%s\n", c.Name, c.Help)
	}
	w.Flush()
	bytes, _ := ioutil.ReadAll(&buf)
	return fmt.Sprintf(help, id, id, string(bytes))
}
