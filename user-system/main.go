package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	log "util/logger"
	"wxappsrv"
)

// var logln = util.Log
// var logf = util.Logf

const (
	DefaultPort = wxappsrv.DefaultPort
)

var (
	port string
	// ...
)

func init() {
	flag.StringVar(&port, "p", DefaultPort, "The PORT to be listened by wxapp.")
}

func main() {
	flag.Parse()
	// TODO: validate port ?

	wxappsrv.LoadAll()
	defer wxappsrv.SaveAll()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Infof("Signal %v", <-c)
		wxappsrv.SaveAll()
		os.Exit(0)
	}()

	err := wxappsrv.Listen(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
