package argparse

import (
	"strings"
)

type Flags struct {
	Verbose bool
}

func CollectFlags(args []string) Flags {
	collected := Flags{}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		if arg == "--verbose" {
			collected.Verbose = true
		}
	}

	return collected
}
