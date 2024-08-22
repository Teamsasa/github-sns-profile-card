package server

import (
	"net/http"
	"path/filepath"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.defaultHandler)
	mux.HandleFunc("/svg", s.SVGHandler)

	assetsDir := http.Dir(filepath.Join("internal", "server", "assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsHandler(assetsDir)))

	return mux
}
