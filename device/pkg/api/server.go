package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server is the struct used to set up the device API
type Server struct {
	router *mux.Router
}

// NewServer creates and returns the reference to a new Server struct
func NewServer(router *mux.Router) *Server {
	s := &Server{
		router: router,
	}
	return s
}

// ListenAndServe makes the server router listen so that the API endpoints are available
func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(":55555", s.router))
}
