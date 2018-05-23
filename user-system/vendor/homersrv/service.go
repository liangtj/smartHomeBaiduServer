package homersrv

import (
	"bytes"
	errors "convention/homererror"
	"encoding/json"
	"entity"
	"fmt"
	"math/rand"
	"model"
	"net/http"
	"strings"
	"time"
	"util"
	log "util/logger"

	"github.com/gin-gonic/gin"
	muxx "github.com/gorilla/mux"

	// "github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gorilla/sessions"
)

var (
	secret = []byte("something-very-secret")
	// SessionStore = cookie.NewStore(key)
	SessionStore = sessions.NewCookieStore(secret)
	sessionName  = "homer-user"
)

type Username = entity.Username
type Auth = entity.Auth

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

	info.Name = username
	info.Auth = password
	info.Mail = email
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

/*
var logInHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "POST")

	// var uInfoRaw UserInfoRaw
	var uInfoRaw struct {
		UserID       string `json:"userId"`
		UserPassword string `json:"userPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		RespondErrorDecoding(w, err)
		return
	}

	userid := Username(uInfoRaw.UserID)
	uInfo, err := QueryAccountByUsername(userid)
	if err != nil {
		RespondError(w, err)
		return
	}

	// LogIn(userid, authTrial)
	authTrial := Auth(uInfoRaw.UserPassword)
	if !uInfo.Auth.Verify(authTrial) {
		RespondError(w, errors.ErrFailedAuth)
	} else {
		maxAge := 10 * time.Minute
		expires := time.Now().Add(maxAge)
		cookie := http.Cookie{
			Name:  "homer-user",
			Value: "",

			// Path:
			// Domain:
			Expires: expires,

			MaxAge: int(maxAge), // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
			// Secure: true,
			// HttpOnly: true,
		}

		res := struct {
			StateCode int `json:"stateCode"`
		}{
			StateCode: 1, // success
		}

		sInfo := entity.SessionInfo{
			ExpiredAt: expires,
			User:      uInfo,
		}
		if err := CreateSession(&sInfo); err != nil {

			// ... for matching need ...
			// RespondError(w, err)
			res.StateCode = 0
			RespondJSON(w, StatusCodeCorrespondingToAgendaError[err], res)

			return
		}
		cookie.Value = string(sInfo.Token)
		http.SetCookie(w, &cookie)

		RespondJSON(w, http.StatusCreated, res)
	}
} */
var logOutHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "DELETE")

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	sInfo, err := Authorize(rInfo.Token)
	if err != nil {
		RespondError(w, err)
		return
	}

	if err := DeleteSession(&sInfo); err != nil {
		RespondError(w, err)
		return
	}

	// RespondJSON(w, http.StatusNoContent)
	// RespondError(w, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}

var getUserByIDHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	if _, err := Authorize(rInfo.Token); err != nil {
		RespondError(w, err)
		return
	}

	if us := muxx.Vars(r)["identifier"]; len(us) > 0 { // FIXME: used muxx
		// if us := r.URL.Query()["username"]; len(us) > 0 {
		// name := Username(us[0])
		name := Username(us)
		uInfo, err := QueryAccountByUsername(name)
		if err != nil {
			RespondError(w, err)
			return
		}

		res := ResponseUserInfoPublic(uInfo.UserInfoPublic)
		RespondJSON(w, http.StatusOK, res)
	}
}
var deleteUserByIDHandler = func(w http.ResponseWriter, r *http.Request) { // Method: "DELETE"
}

var getUsersHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	var rInfo RequestJSON
	if err := json.NewDecoder(r.Body).Decode(&rInfo); err != nil {
		RespondErrorDecoding(w, err)
		return
	}

	if _, err := Authorize(rInfo.Token); err != nil {
		RespondError(w, err)
		return
	}

	// uInfos := QueryAccountAll()
	if uInfos, err := model.UserInfoService.FindAll(); err != nil {
		RespondError(w, err)
	} else {
		res := make([]entity.UserInfoPublic, 0, len(uInfos))
		for _, u := range uInfos {
			res = append(res, u.UserInfoPublic)
		}
		RespondJSON(w, http.StatusOK, res)
	}
}

var registerUserHandler = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "POST")

	var uInfoRaw UserInfoRaw
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		RespondError(w, http.StatusBadRequest, err.Error(), "decode error for elements POST-ed")
		return
	}

	uInfo := MakeUserInfo(
		Username(uInfoRaw.Name),
		Auth(uInfoRaw.Auth),
		uInfoRaw.Mail,
		uInfoRaw.Phone,
	)
	if err := RegisterUser(uInfo); err != nil {
		RespondError(w, err)
		return
	}

	res := ResponseUserInfoPublic(uInfo.UserInfoPublic)
	RespondJSON(w, http.StatusCreated, res)
}

func init() {
	// TODEL: after gin
	// FIXME: when use `curl` and no-trail-slash url to test, fail to be redirected to with-trail-slash version like when using browser .... whatever mux or muxx
	// when using muxx, seems not redirect sub-tree (like '/users/a' --> '/users/') ...
	// mux := mux.NewServeMux()

	router := gin.Default()
	api := "/v1"

	// router.Use(sessions.Sessions("mysession", SessionStore))

	// router.POST(api+"/sessions/", gin.WrapF(logInHandler))
	router.DELETE(api+"/session", gin.WrapF(logOutHandler))

	// Group User
	router.GET(api+"/user/{identifier}", gin.WrapF(getUserByIDHandler))
	router.DELETE(api+"/user/{identifier}", gin.WrapF(deleteUserByIDHandler))
	/* 	router.GET(api+"/user/{identifier}/meetings", gin.WrapF(getMeetingsForUserHandler))
	router.DELETE(api+"/user/{identifier}/meetings", gin.WrapF(deleteMeetingsForUserHandler))
	*/

	// Group Users
	router.GET(api+"/users/", gin.WrapF(getUsersHandler))
	// router.POST(api+"/users/", gin.WrapF(registerUserHandler))

	/* 	// Group Meeting
		router.GET(api+"/meetings/{identifier}", gin.WrapF(getMeetingByIDHandler))
		router.DELETE(api+"/meetings/{identifier}", gin.WrapF(deleteMeetingByIDHandler))
		router.PATCH(api+"/meetings/{identifier}", gin.WrapF(modifyMeetingByIDHandler))

		// Group Meetings
		router.GET(api+"/meetings/", gin.WrapF(getMeetingByIntervalHandler))
	    router.POST(api+"/meetings/", gin.WrapF(sponsorMeetingHandler))
	*/

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

	router.GET(api+"/users/:identifier", func(c *gin.Context) {
		identifier := c.Param("identifier")
		retrieveUserInfoByName(c, Username(identifier))
	})
	router.PATCH(api+"/users/:identifier", func(c *gin.Context) {
		identifier := c.Param("identifier")
		modifyUserInfoByName(c, Username(identifier)) // TODO: AuthRequired
	})
	router.DELETE(api+"/users/:identifier", func(c *gin.Context) {
		identifier := c.Param("identifier")
		deleteUserByName(c, Username(identifier)) // TODO: AuthRequired
	})

	router.GET(api+"/is-new-user", gin.WrapF(isNewUserGet)) // TODEL: Conflict
	router.POST(api+"/isNewUser", isNewUser)

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
		UserID       string `json:"username"`
		UserPassword string `json:"password"`
	}

	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userid := Username(uInfoRaw.UserID)
	uInfo, err := QueryAccountByUsername(userid)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	// LogIn(userid, authTrial)
	authTrial := Auth(uInfoRaw.UserPassword)
	if !uInfo.Auth.Verify(authTrial) {
		err := errors.ErrFailedAuth
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
	} else {
		sess, _ := SessionStore.Get(c.Request, sessionName)

		sess.Values[42] = rand.Uint32()

		maxAge := 10 * time.Minute
		sess.Options.MaxAge = int(maxAge)

		sess.Values["authenticated"] = true
		if err := sess.Save(c.Request, c.Writer); err != nil {
			internalError(c)
			return
		}

		res := struct {
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}{
			Msg: "OK",
			Data: Object{
				"username": userid,
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
var retrieveUserInfoByName = func(c *gin.Context, username Username) {
	uInfo, err := QueryAccountByUsername(username)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"username":       username.String(),
			"homeAssitantIP": uInfo.HomeAssitantAddr,
		},
	}
	c.JSON(200, res)
}
var modifyUserInfoByName = func(c *gin.Context, username Username) {
	var uInfoRaw struct {
		UserID       string `json:"username"`
		UserPassword string `json:"newpassword"`
	}

	uInfo, err := QueryAccountByUsername(username)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	// TODO: modify !

	res := ResponseBody{
		Msg: "OK",
		Data: Object{
			"username": uInfoRaw.UserID,
			"password": uInfoRaw.UserPassword,
		},
	}
	c.JSON(200, res)
}
var deleteUserByName = func(c *gin.Context, username Username) {
	var uInfoRaw struct {
		UserPassword string `json:"password"`
	}
	if err := c.ShouldBind(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	uInfo, err := QueryAccountByUsername(username)
	if err != nil {
		c.AbortWithError(StatusCodeCorrespondingToAgendaError[err], err)
		return
	}

	// TODO: validate && delete !

	res := ResponseBody{
		Msg:  "OK",
		Data: nil,
	}
	c.JSON(200, res)
}

var isNewUserGet = func(w http.ResponseWriter, r *http.Request) {
	util.PanicIf(r.Method != "GET")

	var uInfoRaw struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&uInfoRaw); err != nil {
		// NOTE: maybe should not expose `err` ?
		RespondError(w, http.StatusBadRequest, err.Error(), "decode error for elements GET-ed")
		return
	}

	res := struct {
		IsNewUser bool `json:"isNewUser"`
	}{
		IsNewUser: rand.Float32() < 0.5,
	}
	RespondJSON(w, http.StatusCreated, res)
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
