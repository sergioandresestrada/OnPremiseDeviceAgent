package server

import (
	objstorage "backend/pkg/obj_storage"
	"backend/pkg/queue"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	queue       queue.Queue
	obj_storage objstorage.Obj_storage
	router      *mux.Router
	serverURL   string
}

func NewServer(queue queue.Queue, obj_storage objstorage.Obj_storage, router *mux.Router) *Server {
	url, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		panic("Environment variable SERVER_URL does not exist")
	}

	s := &Server{
		router:      router,
		queue:       queue,
		obj_storage: obj_storage,
		serverURL:   url}
	return s
}

func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(":12345", s.router))
}
