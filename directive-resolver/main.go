package main

import (
	"asr"
	"auth"
	"io/ioutil"
	"nlp/lexer"
	"os"
	"path/filepath"
	"strings"
	"util"
)

var (
	log  = util.Log
	logf = util.Logf
)

// TODO: use lib to detect filetype ?
func filetype(filename string) string {
	t := filepath.Ext(filename)
	if len(t) > 0 && strings.HasPrefix(t, ".") {
		return t[1:]
	}
	return ""
}

func test() {
	filename := "./sample/8k.wav"
	speechFile, err := os.Open(filename)
	util.PanicIf(err)
	defer speechFile.Close()

	speech, err := ioutil.ReadAll(speechFile)
	util.PanicIf(err)

	// token := auth.YuyinBaiduToken()
	token := auth.NLPBaiduToken()

	t, e := asr.SpeechToText(token, asr.Speech{
		Format:  filetype(filename),
		Rate:    8000,
		Channel: 1,
		Speech:  speech,
	})
	util.PanicIf(e)
	log(t)

	// token = auth.NLPBaiduToken()
	lexer.Lexer(token, lexer.Request{
		Text: t,
	})
	util.PanicIf(e)
}

func main() {
	test()
}
