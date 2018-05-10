package util

import (
	"log"
	"smarthome/util"
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

var (
	PanicIf = util.PanicIf
)
