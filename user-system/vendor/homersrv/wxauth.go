package homersrv

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	URL = "https://api.weixin.qq.com/sns/jscode2session"
)

const (
	homerAppID         = "wxa74f1ddc110ca088"
	homerAppSecret     = "c0c5cb53426471ce5c6620d0f439c2c9"
	wxappAuthGrantType = "authorization_code"
)

type WxauthRequest struct {
	AppID     string `json:"appid"`
	Secret    string `json:"secret"`
	JSCode    string `json:"js_code"`
	GrantType string `json:"grant_type"`
}
type WxauthResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`

	UnionID string `json:"unionid"`

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func WxLoginTokenAuth(code string) (WxauthResponse, error) {

	// b := new(bytes.Buffer)
	// if err := json.NewEncoder(b).Encode(wr); err != nil {
	// 	return WxauthResponse{}, err
	// }
	// req, err := http.NewRequest("GET", URL, b)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return WxauthResponse{}, err
	}
	// req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("appid", homerAppID)
	q.Add("secret", homerAppSecret)
	q.Add("js_code", code)
	q.Add("grant_type", wxappAuthGrantType)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WxauthResponse{}, err
	}
	defer resp.Body.Close()

	type Result = WxauthResponse
	var ret WxauthResponse

	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return Result{}, err
	}
	// logf("\n - WxLoginTokenAuth resp: (%+v)\n", ret)
	if ret.ErrCode != 0 && ret.ErrMsg != "" {
		return Result{}, fmt.Errorf(" - error, errcode: %v, errmsg: %v", ret.ErrCode, ret.ErrMsg)
	}
	if ret.UnionID == "" {
		return ret, fmt.Errorf(" - error, 不满足UnionID下发条件的情况")
	}
	return ret, nil
}
