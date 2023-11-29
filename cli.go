package main

import (
	"strings"
)

type Flags struct {
	AddPath   bool
	NoCache   bool
	Verbose   bool
	RawOutput bool
}

// Command definition for CLI subcommands.
type Command struct {
	Label       string
	Handler     func([]string, Flags, State) error
	Usage       string
	Description string
}

// Represents a CLI invocation.
// Must be initialized with commands via AddCommand before running
// with Run.
type CLI struct {
	// Commands in enumeration order.
	OrderedCommands []string
	// Command metadata entries.
	Commands map[string]Command
	Metadata map[string]string
}

// Registers a command.
// This specifies a label that is used to route the user input to
// the right command, a handler that is called when the label is used,
// and usage/description details that get included in autogenerated help messaging.
func (c CLI) AddCommand(label string, handler func([]string, Flags, State) error, usage string, description string) CLI {
	if c.Commands == nil {
		c.Commands = map[string]Command{}
		c.OrderedCommands = []string{}
	}

	c.OrderedCommands = append(c.OrderedCommands, label)
	c.Commands[label] = Command{Label: label, Handler: handler, Usage: usage, Description: description}

	return c
}

// Executes one of the registered commands if any match the provided
// user arguments.
func (c CLI) Run(args []string, currentState State) error {
	if len(args) == 0 {
		c.Help()
		return nil
	}

	command := args[0]

	if command == "help" {
		c.Help()
		return nil
	}

	flags := collectFlags(args)
	return c.Commands[command].Handler(args, flags, currentState)
}

// Prints autogenerated help documentation specifying command usage
// and descriptions based on registered commands (see: AddCommand).
func (c CLI) Help() {
	InfoLogger.Printf("v: A simple version manager. (v%s)\n---", Version)
	for _, commandLabel := range c.OrderedCommands {
		command := c.Commands[commandLabel]
		InfoLogger.Printf("\033[1m%-30s\033[0m%s\n", command.Usage, command.Description)
	}
}

// Traverses input arguments and extracts flags of
// the form --<flag-label>.
func collectFlags(args []string) Flags {
	collected := Flags{}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		switch arg {
		case "--verbose":
			collected.Verbose = true
		case "--no-cache":
			collected.NoCache = true
		case "--add-path":
			collected.AddPath = true
		case "--raw":
			collected.RawOutput = true
		}
	}

	return collected
}
