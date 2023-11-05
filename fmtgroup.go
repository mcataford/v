package main

import (
	"fmt"
)

func StartFmtGroup(label string) func(string) {
	fmt.Printf("\033[1m%s\033[0m\n", label)

	return func(message string) {
		fmt.Printf("    %s\n", message)
	}
}
