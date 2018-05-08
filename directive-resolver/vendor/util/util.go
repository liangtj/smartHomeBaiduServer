package util

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

var (
	Log  = log.Println
	Logf = log.Printf

// Log  = func(args ...interface{}) {}
// Logf = func(args ...interface{}) {}
)
