package api

func (s *Server) Routes() {
	s.router.HandleFunc("/jobs", s.Jobs).Methods("GET")
	s.router.HandleFunc("/identification", s.Identification).Methods("GET")
	s.router.HandleFunc("/job", s.ReceiveJob).Methods("POST")
	s.router.HandleFunc("/heartbeat", s.Heartbeat).Methods("POST")
}
