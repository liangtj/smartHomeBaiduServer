package logger

import (
	"config"
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	logWriter := os.Stderr
	if !config.LogToConsoleMode() {
		flog, err := os.OpenFile(config.LogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			log.Panic(err)
		}
		logWriter = flog
	}
	// logWriter := io.MultiWriter(flog, os.Stderr)

	Logger = log.New(logWriter, "cloudgo: ", log.LstdFlags|log.Lshortfile)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func Info(v ...interface{}) {
	Logger.SetPrefix("[info]")
	Logger.Output(2, fmt.Sprint(v...))
}
func Infof(format string, v ...interface{}) {
	Logger.SetPrefix("[info]")
	Logger.Output(2, fmt.Sprintf(format, v...))
}
func Infoln(v ...interface{}) {
	Logger.SetPrefix("[info]")
	Logger.Output(2, fmt.Sprintln(v...))
}

func Warning(v ...interface{}) {
	Logger.SetPrefix("[warning]")
	Logger.Output(2, fmt.Sprint(v...))
}
func Warningf(format string, v ...interface{}) {
	Logger.SetPrefix("[warning]")
	Logger.Output(2, fmt.Sprintf(format, v...))
}
func Warningln(v ...interface{}) {
	Logger.SetPrefix("[warning]")
	Logger.Output(2, fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Output(2, fmt.Sprint(v...))
}
func Errorf(format string, v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Output(2, fmt.Sprintf(format, v...))
}
func Errorln(v ...interface{}) {
	Logger.SetPrefix("[error]")
	Logger.Output(2, fmt.Sprintln(v...))
}

var (
	Print   = Info
	Printf  = Infof
	Println = Infoln
)

func Fatal(v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Output(2, fmt.Sprint(v...))

	log.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Output(2, fmt.Sprintf(format, v...))

	log.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func Fatalln(v ...interface{}) {
	Logger.SetPrefix("[fatal]")
	Logger.Output(2, fmt.Sprintln(v...))

	log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Panic(v ...interface{}) {
	Logger.SetPrefix("[panic]")
	s := fmt.Sprint(v...)
	Logger.Output(2, s)
	log.Output(2, s)
	panic(s)
}
func Panicf(format string, v ...interface{}) {
	Logger.SetPrefix("[panic]")
	s := fmt.Sprintf(format, v...)
	Logger.Output(2, s)
	log.Output(2, s)
	panic(s)
}
func Panicln(v ...interface{}) {
	Logger.SetPrefix("[panic]")
	s := fmt.Sprintln(v...)
	Logger.Output(2, s)
	log.Output(2, s)
	panic(s)
}

// var (
// 	Print   = log.Print
// 	Printf  = log.Printf
// 	Println = log.Println
// 	Fatal   = log.Fatal
// 	Fatalf  = log.Fatalf
// 	Fatalln = log.Fatalln
// 	Panic   = log.Panic
// 	Panicf  = log.Panicf
// 	Panicln = log.Panicln

// 	// TODO:
// 	Warning   = log.Print
// 	Warningf  = log.Printf
// 	Warningln = log.Println
// 	Error     = log.Print
// 	Errorf    = log.Printf
// 	Errorln   = log.Println
// )
