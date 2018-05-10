package util

import (
	"log"
	"smarthome/util"
)

// var (
// 	// Logger = logger.Logger

// 	Print   = Logger.Print
// 	Printf  = Logger.Printf
// 	Println = Logger.Println
// 	Fatal   = Logger.Fatal
// 	Fatalf  = Logger.Fatalf
// 	Fatalln = Logger.Fatalln
// 	Panic   = Logger.Panic
// 	Panicf  = Logger.Panicf
// 	Panicln = Logger.Panicln

// 	// TODO:
// 	Warning   = Logger.Print
// 	Warningf  = Logger.Printf
// 	Warningln = Logger.Println
// 	Error     = Logger.Print
// 	Errorf    = Logger.Printf
// 	Errorln   = Logger.Println
// )

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
