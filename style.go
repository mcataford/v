package main

const (
	RESET  = "\033[0m"
	BOLD   = "\033[1m"
	YELLOW = "\033[33m"
)

func Yellow(text string) string {
	return YELLOW + text + RESET
}

func Bold(text string) string {
	return BOLD + text + RESET
}
