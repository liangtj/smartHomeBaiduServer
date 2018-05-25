package homersrv

import (
	"bytes"
	"config"
	errors "convention/homererror"
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

	// "github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gorilla/sessions"
)

var (
	storeFilepath = config.WorkingDir() + "session-store.json"
	secret        = []byte("something-very-secret")
	// SessionStore = cookie.NewStore(key)
	SessionStore = sessions.NewFilesystemStore(storeFilepath, secret)
	sessionName  = "homer-user"
)

type Username = entity.Username
type Auth = entity.Secret

// type UserInfo = entity.UserInfoRaw
type UserInfoRaw struct {
	Name  string `json:"username"`
	Auth  string `json:"password"`
	Mail  string `json:"mail"`
	Phone string `json:"phone"`
}

type RequestJSON struct {
	Token entity.Token `json:"token"`
	UserInfoRaw
	// ...
}

type UserInfoPublic = entity.UserInfoPublic
type User = entity.User
type MeetingInfo = entity.MeetingInfo
type Meeting = entity.Meeting
type MeetingTitle = entity.MeetingTitle

func MakeUserInfo(username Username, password Auth, email, phone string) entity.UserInfo {
	info := entity.UserInfo{}

	info.ID = username
	info.Secret = password
	info.HomeAssistantAddr = email
	info.Phone = phone

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
	homerUserSys struct {
		// *server.Server
		Server *gin.Engine
	}
)

func init() {
	// TODEL: after gin
	// FIXME: when use `curl` and no-trail-slash url to test, fail to be redirected to with-trail-slash version like when using browser .... whatever mux or muxx
	// when using muxx, seems not redirect sub-tree (like '/users/a' --> '/users/') ...
	// mux := mux.NewServeMux()

	router := gin.Default()
	api := "/v1"

	// router.Use(sessions.Sessions("mysession", SessionStore))

	// ...
	router.GET("/api/test", gin.WrapF(apiTestHandler()))
	router.GET("/unknown/", gin.WrapF(sayDeveloping))
	router.GET("/say/", gin.WrapF(sayhelloName))

	// With gin, should use `StaticFS` to let it work like a FS ;
	// Or, using `Static` would need something like `http.StripPrefix` ...
	router.StaticFS("/static", http.Dir("./asset"))

	// @@binly:

	router.POST(api+"/login", login)

	router.POST(api+"/logout", logout)

	router.POST(api+"/register", register) // TODEL:
	router.POST(api+"/users", register)

	router.GET(api+"/userInfo", retrieveUserInfo)
	router.PATCH(api+"/userInfo", modifyUserInfo) // TODO: AuthRequired

	router.POST(api+"/isNewUser", isNewUser)

	// TODEL: TEST
	router.GET(api+"/test-login", loginTest)

	homerUserSys.Server = router
}

// @@binly:

type Object = map[string]interface{}

type ResponseBody struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var register = func(c *gin.Context) {
	var uInfoRaw struct {
		UserID           string `json:"username"`
		UserPassword     string `json:"password"`
		HomeAssitantAddr string `json:"homeAssitantIP"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userid := Username(uInfoRaw.UserID)
	_, err := QueryAccountByUsername(userid)
	if err == nil {
		c.JSON(409, ResponseBody{
			Msg:  "Conflict",
			Data: nil,
		})
		return
	}
	if err != errors.ErrUserNotFound && err != errors.ErrUserNotRegistered {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	HAAddr := uInfoRaw.HomeAssitantAddr

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"username":       userid,
			"homeAssitantIP": HAAddr,
		},
	}
	c.JSON(200 /* http.StatusCreated */, res)
}

// var login = logInHandler
var login = func(c *gin.Context) {
	var uInfoRaw struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	wxResp, err := WxLoginTokenAuth(uInfoRaw.Code)

	if err != nil {
		err := errors.ErrFailedAuth
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
	} else {

		// if need, use whatever a session-mgr ...

		u := c.Request.URL.String()
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			internalError(c)
			return
		}
		// context.Set(req, registryKey, ...)

		// sessToken3rd := "...3rd_session..."
		sessToken3rd := entity.TokenGen(128)
		req.Header.Set("Wx-3rd-Session-Token", sessToken3rd.String())

		sess, _ := SessionStore.Get(req, sessionName)

		sess.Values[42] = rand.Uint32()
		sess.Values["3rd_session"] = sessToken3rd
		sess.Values["Wx-Open-ID"] = wxResp.OpenID
		sess.Values["Wx-Session-Key"] = wxResp.SessionKey

		maxAge := 10 * time.Minute
		sess.Options.MaxAge = int(maxAge)

		sess.Values["authenticated"] = true
		if err := sess.Save(req /* c.Request */, c.Writer); err != nil {
			internalError(c)
			return
		}

		res := struct {
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}{
			Msg: "OK",
			Data: Object{
				"3rd_session": sess.Values["3rd_session"],
			},
		}
		c.JSON(200 /* http.StatusCreated */, res)
	}
}

var logout = func(c *gin.Context) {
	var uInfoRaw struct {
		UserID string `json:"username"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
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

}

// FIXME: AuthRequired ?
var retrieveUserInfo = func(c *gin.Context) {
	sess, err := SessionStore.Get(c.Request, sessionName)
	if err != nil {
		internalError(c)
	}

	id := sess.Values["Wx-Open-ID"].(Username)
	uInfo, err := QueryAccountByUsername(id)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
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
	var uInfoRaw struct {
		HomeAssistantAddr string `json:"homeAssistantIP"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sess, err := SessionStore.Get(c.Request, sessionName)
	if err != nil {
		internalError(c)
	}

	id := sess.Values["Wx-Open-ID"].(Username)
	uInfo, err := QueryAccountByUsername(id)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	// TODO: modify !
	uInfo.HomeAssistantAddr = uInfoRaw.HomeAssistantAddr
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
var isNewUser = func(c *gin.Context) {
	var uInfoRaw struct {
		UserID       string `json:"username"`
		UserPassword string `json:"password"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userid := Username(uInfoRaw.UserID)
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
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	c.JSON(403 /* ... */, ResponseBody{
		Msg:  "Forbidden",
		Data: nil,
	})
}

// TODEL: working with `gin-contrib/sessions`
/* func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You'd normally redirect to login page
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		} else {
			// Continue down the chain to handler etc
			c.Next()
		}
	}
} */
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := SessionStore.Get(c.Request, sessionName)
		if err != nil {
			internalError(c)
		}
		if sess == nil || !sess.Values["authenticated"].(bool) {
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

// @@binly:

func Listen(addr string) error {
	if addr == "" {
		addr = DefaultPort
	}
	return homerUserSys.Server.Run(addr)
	// return homerUserSys.Listen(addr)
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
