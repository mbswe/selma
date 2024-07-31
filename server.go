package selma

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Config holds the server configuration

// Server holds settings for the HTTP server
type Server struct {
	Addr   string
	Router *Router
	Config *Config
}

// NewServer creates a new server with provided settings
func NewServer(configPath string, router *Router) *Server {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	addr := ":" + strconv.Itoa(config.ServerPort)
	return &Server{
		Addr:   addr,
		Router: router,
		Config: config,
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe() {
	log.Printf("Server listening on %s", s.Addr)
	err := http.ListenAndServe(s.Addr, s.Router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
