package wxappsrv

import (
	"encoding/json"
	"log"
	"math/rand"
	"model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var loginTest = func(c *gin.Context) {
	log.Print(".........>")

	// now := time.Now()
	// wxResp := WxauthResponse{
	// 	OpenID:     now.String(),
	// 	SessionKey: string(time.Now().Sub(now)),
	// 	UnionID:    string(time.Now().Sub(now)),
	// }
	now := strconv.FormatUint(rand.Uint64()%10, 10)
	wxResp := WxauthResponse{
		OpenID:     now,
		SessionKey: now,
		UnionID:    now,
	}

	log.Printf(" ---\n wxResp:(%+v)\n\n", wxResp)

	// sessToken3rd := "...3rd_session..."
	// sessToken3rd := entity.TokenGen(128).String()
	sessToken3rd := "XXX" + strconv.FormatUint(rand.Uint64()%10, 10)
	maxAge := 10 * time.Minute
	expiry := time.Now().Add(maxAge)
	sessData, err := json.Marshal(wxResp)
	if err != nil || (sessionStore.Save(sessToken3rd, sessData, expiry) != nil) {
		internalError(c)
		return
	}
	log.Printf(" --- sessToken:(%+v), sessData:(%+v)\n", sessToken3rd, wxResp)

	// log.Printf(" ---\n sess:(%+v),\n len:(%v)\n\n", sess, len(SessionStore.Codecs))

	id := wxResp.OpenID
	if _, err := QueryAccountByUsername(Username(id)); err != nil {
		log.Printf("err: %v", err)
		if err == gorm.ErrRecordNotFound {
			uInfo := MakeUserInfo(Username(id), "")
			if e := model.UserInfoService.Create(&uInfo); e != nil {
				c.AbortWithError(StatusCodeCorrespondingToWxappError[e], e)
				return
			}
		} else {
			c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
			return
		}
	}

	res := struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Msg: "OK",
		Data: Object{
			"3rd_session": sessToken3rd,
		},
	}
	c.JSON(200 /* http.StatusCreated */, res)

}

var retrieveUserInfoTest = func(c *gin.Context) {
	sessToken, sessData, err := retrieveWxappSession(c)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
		return
	}

	// id := sessData.OpenID
	// uInfo, err := QueryAccountByUsername(Username(id))
	// if err != nil {
	// 	c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
	// 	return
	// }

	log.Printf(" --- sessToken:(%+v), sessData:(%+v)\n", sessToken, sessData)

	res := ResponseBody{
		Msg:  "OK",
		Data: Object{
		// "homeAssistantIP": uInfo.HomeAssistantAddr,
		},
	}
	c.JSON(200, res)
}
