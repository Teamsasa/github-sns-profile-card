package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.DefaultHandler)
	mux.HandleFunc("/svg", s.SVGHandler)

	return mux
}

func (s *Server) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
