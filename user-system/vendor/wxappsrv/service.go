package wxappsrv

import (
	"bytes"
	errors "convention/errors"
	"encoding/json"
	"entity"
	"fmt"
	"math/rand"
	"model"
	"net/http"
	"strings"
	"time"
	log "util/logger"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/alexedwards/scs/stores/memstore"
)

var (
	sessionName = "wxapp-user"

	sessionStoreCleanupInterval = 10 * time.Second
	sessionStore                = memstore.New(sessionStoreCleanupInterval)
)

type Username = entity.UserIdentifier
type Auth = entity.Secret

type UserInfoPublic = entity.UserInfoPublic
type User = entity.User

func MakeUserInfo(id Username, HAAddr string) entity.UserInfo {
	info := entity.UserInfo{}

	info.ID = id
	// info.Secret = password
	info.HomeAssistantAddr = HAAddr
	// info.Phone = phone

	return info
}

func LoadAll() {
	model.Load()
}
func SaveAll() {
	if err := model.Save(); err != nil {
		log.Error(err)
	}
}

// Server ...

const (
	DefaultPort = "8080"
)

var (
	wxappUserSys struct {
		Server *gin.Engine
	}
)

func init() {
	router := gin.Default()
	api := "/v1"

	// ...
	router.GET("/api/test", gin.WrapF(apiTestHandler()))
	router.GET("/unknown/", gin.WrapF(sayDeveloping))
	router.GET("/say/", gin.WrapF(sayhelloName))

	// With gin, should use `StaticFS` to let it work like a FS ;
	// Or, using `Static` would need something like `http.StripPrefix` ...
	router.StaticFS("/static", http.Dir("./asset"))

	router.POST(api+"/login", login)
	// router.POST(api+"/logout", logout)

	authorized := router.Group("/", AuthRequired())
	authorized.GET(api+"/userInfo", retrieveUserInfo)
	authorized.PATCH(api+"/userInfo", modifyUserInfo)

	// router.POST(api+"/isNewUser", isNewUser)

	// TODEL: TEST
	router.GET(api+"/test-login", loginTest)
	authorized.GET(api+"/test-ru", retrieveUserInfoTest)

	wxappUserSys.Server = router
}

type Object = map[string]interface{}

type ResponseBody struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var WxappSessionToken = func(c *gin.Context) string {
	t := c.Request.Header.Get("Wx-Session-Token")
	return t
}
var setWxappSessionToken = func(c *gin.Context, t string) {
	c.Request.Header.Set("Wx-Session-Token", t)
}

var login = func(c *gin.Context) {
	var reqData struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBind(&reqData); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	wxResp, err := WxLoginTokenAuth(reqData.Code)

	if err != nil {
		err := errors.ErrFailedAuth
		c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
	} else {
		sessToken := entity.TokenGen(128).String()
		maxAge := 10 * time.Minute
		expiry := time.Now().Add(maxAge)
		sessData, err := json.Marshal(wxResp)
		if err != nil || (sessionStore.Save(sessToken, sessData, expiry) != nil) {
			internalError(c)
			return
		}

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
}

/* var logout = func(c *gin.Context) {
	var reqData struct {
		UserID string `json:"username"`
	}

	if err := c.ShouldBind(&reqData); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sess, err := SessionStore.Get(c.Request, sessionName)
	if err != nil {
		internalError(c)
	}

	sess.Values["authenticated"] = false
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request, c.Writer); err != nil {
		internalError(c)
		return
	}

	res := ResponseBody{
		Msg:  "OK",
		Data: nil,
	}
	c.JSON(200, res)

} */

// TODO: let be middleware ?
var retrieveWxappSession = func(c *gin.Context) (string, WxauthResponse, error) {
	t := WxappSessionToken(c)
	sessB, _, err := sessionStore.Find(t)
	if err != nil {
		// internalError(c) TODO: if middleware ...
		return t, WxauthResponse{}, err
	}
	var sessData WxauthResponse
	if err := json.Unmarshal(sessB, &sessData); err != nil {
		// internalError(c) TODO: if middleware ...
		return t, sessData, err
	}
	return t, sessData, nil
}

var retrieveUserInfo = func(c *gin.Context) {
	/* t := WxappSessionToken(c)
		sessB, _, err := sessionStoreMem.Find(t)
		if err != nil {
			internalError(c)
	    }
	    var sessData WxauthResponse
	    err := json.Unmarshal(sessB, &sessData)
	    if err != nil {
	        internalError(c)
	        return
	    } */
	_, sessData, err := retrieveWxappSession(c)
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

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"homeAssistantIP": uInfo.HomeAssistantAddr,
		},
	}
	c.JSON(200, res)
}
var modifyUserInfo = func(c *gin.Context) {
	var reqData struct {
		HomeAssistantAddr string `json:"homeAssistantIP"`
	}

	if err := c.ShouldBind(&reqData); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, sessData, err := retrieveWxappSession(c)
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

	// TODO: modify !
	uInfo.HomeAssistantAddr = reqData.HomeAssistantAddr
	if err := model.UserInfoService.Save(&uInfo); err != nil {
		res := ResponseBody{
			Msg:  "Forbidden",
			Data: nil,
		}
		log.Error(err)
		c.JSON(403, res)
	} else {
		res := ResponseBody{
			Msg: "OK",
			Data: Object{
				"homeAssistantAddr": uInfo.HomeAssistantAddr,
			},
		}
		c.JSON(200, res)
	}
}

/* var isNewUser = func(c *gin.Context) {
	var reqData struct {
		UserID       string `json:"username"`
		UserPassword string `json:"password"`
	}

	if err := c.ShouldBind(&reqData); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userid := Username(reqData.UserID)
	_, err := QueryAccountByUsername(userid)
	if err == nil {
		res := ResponseBody{
			Msg: "OK",
			Data: Object{
				"username": userid,
			},
		}
		c.JSON(200, res)
	}
	if err != errors.ErrUserNotFound && err != errors.ErrUserNotRegistered {
		c.AbortWithError(StatusCodeCorrespondingToWxappError[err], err)
		return
	}

	c.JSON(403, ResponseBody{
		Msg:  "Forbidden",
		Data: nil,
	}) // ...
} */

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := WxappSessionToken(c)
		_, exist, err := sessionStore.Find(t)
		if err != nil {
			internalError(c)
		} else if !exist {
			c.JSON(http.StatusUnauthorized, ResponseBody{
				Msg:  "Unauthorized",
				Data: nil,
			})
		} else {
			c.Next()
		}
	}
}

var internalError = func(c *gin.Context) {
	// http.Error(w, "", http.StatusInternalServerError)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func Listen(addr string) error {
	if addr == "" {
		addr = DefaultPort
	}
	return wxappUserSys.Server.Run(addr)
	// return wxappUserSys.Listen(addr)
}

// detail handlers, etc ... ----------------------------------------------------------------

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	segments := strings.Split(r.URL.Path, "/")
	name := segments[len(segments)-1]
	fmt.Fprintf(w, "Hello %v!\n", name)
}

func sayDeveloping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)

	fmt.Fprintf(w, "Developing!\n")
	fmt.Fprintf(w, "Now NotImplemented!\n")
}

func apiTestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "9527", Content: "Hello from Go!\n"}

		// json.NewEncoder(w).Encode(res)
		j, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rand.Seed(time.Now().UnixNano())
		prettyPrint := rand.Float32() < 0.5
		if prettyPrint {
			var out bytes.Buffer
			json.Indent(&out, j, "", "\t")
			j = out.Bytes()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}
}
