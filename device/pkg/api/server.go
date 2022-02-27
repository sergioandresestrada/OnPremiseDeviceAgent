package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(router *mux.Router) *Server {
	s := &Server{
		router: router,
	}
	return s
}

func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(":55555", s.router))
}
