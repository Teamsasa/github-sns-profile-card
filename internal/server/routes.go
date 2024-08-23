package server

import (
	"net/http"
	"path/filepath"
)

func assetsHandler(dir http.Dir) http.Handler {
	return http.FileServer(dir)
}

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.DefaultHandler)
	mux.HandleFunc("/svg", s.SVGHandler)

	assetsDir := http.Dir(filepath.Join("internal", "server", "assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsHandler(assetsDir)))

	return mux
}

func (s *Server) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic for the default handler here
}
