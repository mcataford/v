package main

import "fmt"

const (
	RESET = "\033[0m"
	BOLD  = "\033[1m"
)

func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", BOLD, text, RESET)
}
