package main

import "fmt"

const (
	RESET  = "\033[0m"
	BOLD   = "\033[1m"
	YELLOW = "\033[33m"
)

func Yellow(text string) string {
	return fmt.Sprintf("%s%s%s", YELLOW, text, RESET)
}

func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", BOLD, text, RESET)
}
