package argparse

import (
	"fmt"
	"strings"
	stateManager "v/internal/state"
)

type Command struct {
	Label       string
	Handler     func([]string, Flags, stateManager.State) error
	Usage       string
	Description string
}

type CLI struct {
	Commands map[string]Command
}

func (c CLI) AddCommand(label string, handler func([]string, Flags, stateManager.State) error, usage string, description string) CLI {
	if c.Commands == nil {
		c.Commands = map[string]Command{}
	}

	c.Commands[label] = Command{Label: label, Handler: handler, Usage: usage, Description: description}

	return c
}

func (c CLI) Run(args []string, currentState stateManager.State) error {
	command := args[0]

	if command == "help" {
		c.Help()
		return nil
	}

	flags := CollectFlags(args)
	return c.Commands[command].Handler(args, flags, currentState)
}

func (c CLI) Help() {
	usageStrings := []string{}

	for _, command := range c.Commands {
		usageStrings = append(usageStrings, fmt.Sprintf("\033[1m%-30s\033[0m%s", command.Usage, command.Description))
	}

	helpString := fmt.Sprintf(`v: A simple version manager.
---
%s`, strings.Join(usageStrings, "\n"))

	fmt.Println(helpString)
}
