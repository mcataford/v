package cli

import (
	"strings"
	logger "v/logger"
	state "v/state"
)

type Flags struct {
	AddPath   bool
	NoCache   bool
	Verbose   bool
	RawOutput bool
}

// Represents a CLI invocation.
// Must be initialized with commands via AddCommand before running
// with Run.
type CLI struct {
	Namespaces map[string]Namespace
	Metadata   map[string]string
}

func (c *CLI) AddNamespace(namespace Namespace) *CLI {
	if c.Namespaces == nil {
		c.Namespaces = map[string]Namespace{}
	}

	c.Namespaces[namespace.Label] = namespace

	return c
}

// Executes one of the registered commands if any match the provided
// user arguments.
func (c CLI) Run(args []string, currentState state.State) error {
	if len(args) == 0 {
		c.Help()
		return nil
	}

	action := args[0]

	if action == "help" {
		c.Help()
		return nil
	}

	flags := collectFlags(args)

	namespace, isNamespace := c.Namespaces[action]

	if isNamespace {
		action = args[1]
		return namespace.Commands[action].Handler(args[1:], flags, currentState)
	}

	rootNamespace := c.Namespaces[""]

	return rootNamespace.Commands[action].Handler(args, flags, currentState)
}

// Prints autogenerated help documentation specifying command usage
// and descriptions based on registered commands (see: AddCommand).
func (c CLI) Help() {
	logger.InfoLogger.Printf("v: A simple version manager. (v%s)\n---", c.Metadata["Version"])
	for _, namespace := range c.Namespaces {
		for _, commandLabel := range namespace.ListCommands() {
			command := namespace.Commands[commandLabel]
			logger.InfoLogger.Printf("\033[1m%-30s\033[0m%s\n", command.Usage, command.Description)
		}
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
