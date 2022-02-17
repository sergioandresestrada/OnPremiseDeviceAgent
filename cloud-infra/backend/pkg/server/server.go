package server

import (
	objstorage "backend/pkg/obj_storage"
	"backend/pkg/queue"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	queue       queue.Queue
	obj_storage objstorage.Obj_storage
	router      *mux.Router
}

func NewServer(queue queue.Queue, obj_storage objstorage.Obj_storage, router *mux.Router) *Server {
	s := &Server{
		router:      router,
		queue:       queue,
		obj_storage: obj_storage}
	return s
}

func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(":12345", s.router))
}
