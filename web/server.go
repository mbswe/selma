package web

import (
	"log"
	"net/http"
)

// Server holds settings for the HTTP server
type Server struct {
	Addr   string
	Router *Router
}

// NewServer creates a new server with provided settings
func NewServer(addr string, router *Router) *Server {
	return &Server{
		Addr:   addr,
		Router: router,
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() {
	log.Printf("Server listening on %s", s.Addr)
	http.ListenAndServe(s.Addr, s.Router)
}
