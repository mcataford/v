package logger

import (
	"log"
	"os"
)

var (
	InfoLogger = log.New(os.Stdout, "", 0)
)
