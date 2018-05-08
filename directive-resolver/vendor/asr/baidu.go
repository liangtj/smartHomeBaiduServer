package asr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"util"
)

var (
	log  = util.Log
	logf = util.Logf
)

// ....
const (
	BaiduURL = "http://vop.baidu.com/server_api"
)

// Speech holds the **info** and **content** of speech
// However, later on may be refactored ... (maybe use `anonymous member` ?
type Speech struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	Channel int    `json:"channel"`

	DevPid int `json:"dev_pid,omitempty"`

	Speech []byte `json:"speech"`
}

// Request for ASR
type Request struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	Channel int    `json:"channel"`
	Cuid    string `json:"cuid"`
	Token   string `json:"token"`

	DevPid int `json:"dev_pid,omitempty"`

	// TODO: check **pair** ?
	URL      string `json:"url,omitempty"`
	Callback string `json:"callback,omitempty"`

	Speech string `json:"speech,omitempty"`
	Len    int    `json:"len,omitempty"`
}

// Result : `result` type in baidu-asr's response
type Result []string

// Response for ASR
type Response struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	SN     string `json:"sn"`

	Result Result `json:"result"`
}

// SpeechToText ...
func SpeechToText(token string, speech Speech) (string, error) {
	asrReq := Request{
		Format:  speech.Format,
		Rate:    speech.Rate, // FIXME: Fixed to be 16000 ?
		Channel: 1,
		Cuid:    "homer-test",
		Token:   token,
		// DevPid: 1537,

		Speech: base64.StdEncoding.EncodeToString(speech.Speech),
		Len:    len(speech.Speech),
	}

	ret, err := ASR(asrReq)
	if len(ret) == 0 {
		return "", err
	}
	return ret[0], err
}

// ASR ...
func ASR(asrReq Request) (Result, error) {

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(asrReq); err != nil {
		return Result{}, err
	}

	req, err := http.NewRequest("POST", BaiduURL, b)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	var ret Response
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return Result{}, err
	}
	logf("\n - asr resp: (%+v)\n", ret)
	if ret.ErrNo != 0 || len(ret.Result) == 0 {
		return Result{}, fmt.Errorf(" - error, err_no: %v, err_msg: %v", ret.ErrNo, ret.ErrMsg)
	}
	return ret.Result, nil
}
