package server

import (
	// errors "convention/cloudgoerror"

	"fmt"
	"net/http"
	log "util/logger"
)

type Server struct {
	*http.Server
}

func NewServer() *Server {
	server := new(Server)
	server.Server = &http.Server{}

	return server
}

func (srv *Server) SetHandler(handler http.Handler) error {
	srv.Handler = handler
	return nil
}

func (srv *Server) Listen(addr string) error {
	srv.Addr = addr
	log.Printf("Listtening addr: %v\n", addr)
	fmt.Printf("Listtening addr: %v\n", addr)
	return srv.ListenAndServe()
}
