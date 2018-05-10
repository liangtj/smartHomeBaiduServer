package mux

import "net/http"

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	// TODO:
	mux := new(ServeMux)
	mux.ServeMux = new(http.ServeMux)
	return mux
}
