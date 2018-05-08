package lexer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"util"
)

var (
	log  = util.Log
	logf = util.Logf
)

// ....
const (
	BaiduURL              = "https://aip.baidubce.com/rpc/2.0/nlp/v1/lexer"
	BaiduURLCustomVersion = "https://aip.baidubce.com/rpc/2.0/nlp/v1/lexer_custom"
)

type Request struct {
	Text string `json:"text"` // NOTE: required GEK ?
}

type Response struct {
	Text  string                   `json:"text"`
	Items []map[string]interface{} `json:"items"`
}

// Lexer ...
func Lexer(token string, lexReq Request) (Response, error) {

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(lexReq); err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", BaiduURL, b)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("charset", "UTF-8")
	q.Add("access_token", token)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	var ret Response
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return Response{}, err
	}
	// logf("\n - lexer req: (%+v)\n", req)
	logf("\n - lexer resp: (%+v)\n", ret)
	return ret, nil
}
