package server

import "net/http"

type Server struct {
	addr string
	mux  *http.ServeMux
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) AddHandler(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) AddHandlerFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func NewServer(addr string) *Server {
	mux := &http.ServeMux{}
	return &Server{
		addr: addr,
		mux:  mux,
	}
}
