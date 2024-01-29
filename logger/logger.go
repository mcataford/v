package logger

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "", 0)
	DebugLogger = log.New(ioutil.Discard, "", 0)
)
