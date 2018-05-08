package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"util"
)

var (
	log  = util.Log
	logf = util.Logf
)

// ....
const (
	BaiduURLForToken = "https://openapi.baidu.com/oauth/2.0/token"

	grantType = "client_credentials"
)

var (
	AppInfos = []map[string]string{
		{ // yuyin
			"AppName":   "homer",
			"AppID":     "11199051",
			"APIKey":    "WbYTIVunFdoXRrNB4oo4AoVE",
			"SecretKey": "74d918bcc1e00292608970f7ae63438c",
		},
		{ // nlp
			"AppName":   "homer",
			"AppID":     "11208133",
			"APIKey":    "MGPHeIAcmHaDqB2uEpOeNP2U",
			"SecretKey": "zjKl72Vq7pvEOGpEPs6WUHbkNvjCC5aR",
		},
	}
)

// BaiduOAuth ...
type BaiduOAuth struct {
	clientID     string
	clientSecret string
	// TODO: cache ?
}

type responseOfBaiduOAuth struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// other fields ...

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// NewBaiduOAuth ...
func NewBaiduOAuth(clientID, clientSecret string) *BaiduOAuth {
	return &BaiduOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// GetToken ...
func (oauth *BaiduOAuth) GetToken() (string, error) {
	param := url.Values{}
	param.Add("grant_type", grantType)
	param.Add("client_id", oauth.clientID)
	param.Add("client_secret", oauth.clientSecret)
	resp, err := http.PostForm(BaiduURLForToken, param)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ret responseOfBaiduOAuth
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return "", err
	}
	logf("\n - OAuth resp: (%+v)\n", ret)

	if ret.Error != "" {
		return "", fmt.Errorf(" - error, error: %v, error_description: %v", ret.Error, ret.ErrorDescription)
	}
	return ret.AccessToken, nil
}

// GetBaiduToken ...
func GetBaiduToken(apiKey, secretKey string) string {
	token, err := NewBaiduOAuth(apiKey, secretKey).GetToken()
	util.PanicIf(err) // FIXME:
	return token
}

// YuyinBaiduToken ... return whatever a baidu token TODEL: TEST
func YuyinBaiduToken() string {
	ak := AppInfos[0]["APIKey"]
	sk := AppInfos[0]["SecretKey"]
	return GetBaiduToken(ak, sk)
}

// NLPBaiduToken ... return whatever a baidu token TODEL: TEST
func NLPBaiduToken() string {
	ak := AppInfos[1]["APIKey"]
	sk := AppInfos[1]["SecretKey"]
	return GetBaiduToken(ak, sk)
}
