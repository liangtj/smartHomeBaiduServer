package wxappsrv

import (
	errors "convention/errors"
	"net/http"
)

type HTTPStatusCode = int

// ErrorOrCode can only hold `error` or `HTTPStatusCode` type
type ErrorOrCode = interface{}

var ErrInvalidMethod = errors.New("invalid request method")
var ErrInvalidToken = errors.New("invalid token")
var ErrSessionExpired = errors.New("token fails to be authorized, since sesssion has been expired")

// var ErrDeletedSession = errors.New("deleted session successfully")

var StatusCodeCorrespondingToWxappError = map[error]HTTPStatusCode{
	errors.ErrInvalidUsername: http.StatusBadRequest,
	errors.ErrExistedUser:     http.StatusConflict,
	ErrInvalidMethod:          http.StatusBadRequest,
	errors.ErrFailedAuth:      http.StatusUnauthorized,
	ErrInvalidToken:           http.StatusUnauthorized,
	ErrSessionExpired:         http.StatusUnauthorized,
	// ErrDeletedSession:         http.StatusNoContent,
}
