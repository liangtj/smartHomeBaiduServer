package homersrv

import (
	errors "convention/homererror"
	"encoding/json"
	"entity"
	"fmt"
	"net/http"
	"strings"
	log "util/logger"
)

type HTTPStatusCode = int

// ErrorOrCode can only hold `error` or `HTTPStatusCode` type
type ErrorOrCode = interface{}

var ErrInvalidMethod = errors.New("invalid request method")
var ErrInvalidToken = errors.New("invalid token")
var ErrSessionExpired = errors.New("token fails to be authorized, since sesssion has been expired")

// var ErrDeletedSession = errors.New("deleted session successfully")

var StatusCodeCorrespondingToAgendaError = map[error]HTTPStatusCode{
	errors.ErrInvalidUsername: http.StatusBadRequest,
	errors.ErrExistedUser:     http.StatusConflict,
	ErrInvalidMethod:          http.StatusBadRequest,
	errors.ErrFailedAuth:      http.StatusUnauthorized,
	ErrInvalidToken:           http.StatusUnauthorized,
	ErrSessionExpired:         http.StatusUnauthorized,
	// ErrDeletedSession:         http.StatusNoContent,
}

func RespondError(w http.ResponseWriter, err ErrorOrCode, msg ...string) {
	errString := strings.Join(msg, "\n")
	errCode := http.StatusInternalServerError

	switch e := err.(type) {
	case error:
		errString = e.Error() + "\n\n" + errString
		code, ok := StatusCodeCorrespondingToAgendaError[e]
		if ok {
			errCode = code
		}
	case HTTPStatusCode:
		errCode = e
	default:
		log.Panicf("type `ErrorOrCode` expects `error` or `HTTPStatusCode`, but not %T", e)
	}

	// NOTE: seems that only using `http.Error` to handle simple error is enough ...
	// w.WriteHeader(code)
	// res := ResponseJSON{Error: errString}
	// json.NewEncoder(w).Encode(res)

	http.Error(w, errString, errCode)
}

func RespondErrorDecoding(w http.ResponseWriter, errWhenDecoding error) {
	RespondError(w, http.StatusBadRequest, errWhenDecoding.Error(), "decode error for elements in request")
}

type ResponseJSON struct {
	Error   string      `json:"error"`
	Content interface{} `json:"content"`
}
type ResponseToken struct {
	Token entity.Token `json:"token"`
}
type ResponseUserInfoPublic = entity.UserInfoPublic

func RespondJSON(w http.ResponseWriter, code HTTPStatusCode, res interface{}) { // originally, using `res ResponseJSON`
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

type HTTPMethod = string
type HandlerMap = map[HTTPMethod]http.HandlerFunc

func HandlerMapper(mapping HandlerMap) http.HandlerFunc {
	methods := make([]string, 0, len(mapping))
	for m := range mapping {
		methods = append(methods, m)
	}
	wantedMethods := strings.Join(methods, "/")

	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := mapping[r.Method]
		if ok {
			handler(w, r)
		} else {
			RespondError(w, ErrInvalidMethod, fmt.Sprintf("used method: %v, however, wanted: %v", r.Method, wantedMethods))
			return
		}
	}
}
