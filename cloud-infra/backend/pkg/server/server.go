package server

import (
	"backend/pkg/database"
	objstorage "backend/pkg/obj_storage"
	"backend/pkg/queue"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Server is the struct used to set up the device API.
// It contains a queue and object storage implementation, a rotuer and its own public URL
type Server struct {
	queue      queue.Queue
	objStorage objstorage.ObjStorage
	database   database.Database
	router     *mux.Router
	serverURL  string
}

// NewServer creates and returns the reference to a new Server struct
// It sets the serverURL field to the corresponding Environment variable value, and panics if it not present
func NewServer(queue queue.Queue, objStorage objstorage.ObjStorage, database database.Database, router *mux.Router) *Server {
	url, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		panic("Environment variable SERVER_URL does not exist")
	}

	s := &Server{
		router:     router,
		queue:      queue,
		objStorage: objStorage,
		database:   database,
		serverURL:  url}
	return s
}

// ListenAndServe makes the server router listen so that the API endpoints are available
func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(":12345", s.router))
}
