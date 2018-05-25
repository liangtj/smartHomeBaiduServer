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
	x := "test-" + strconv.FormatUint(rand.Uint64()%5, 10)
	wxResp := WxauthResponse{
		OpenID:     x,
		SessionKey: x,
		UnionID:    x,
	}

	// sessToken := "...3rd_session..."
	// sessToken := entity.TokenGen(128).String()
	sessToken := "XXX" + strconv.FormatUint(rand.Uint64()%5, 10)
	maxAge := 10 * time.Minute
	expiry := time.Now().Add(maxAge)
	sessData, err := json.Marshal(wxResp)
	if err != nil || (sessionStore.Save(sessToken, sessData, expiry) != nil) {
		internalError(c)
		return
	}
	log.Printf(" --- sessToken:(%+v), sessData:(%+v)\n", sessToken, wxResp)

	id := wxResp.OpenID
	if _, err := QueryAccountByUsername(Username(id)); err != nil {
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

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"3rd_session": sessToken,
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

	id := sessData.OpenID
	uInfo, err := QueryAccountByUsername(Username(id))
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
		return
	}

	log.Printf(" --- sessToken:(%+v), sessData:(%+v)\n", sessToken, sessData)

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"homeAssistantIP": uInfo.HomeAssistantAddr,
		},
	}
	c.JSON(200, res)
}
