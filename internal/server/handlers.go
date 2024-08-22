package server

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

// defaultHandler returns a 404 error for all requests. If a route is not found, this handler is called.
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

// assetsHandler serves static files from the assets directory. If the file is not found, it returns a 404 error.
func assetsHandler(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if stat.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func (s *Server) SVGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	iconURL := "/assets/instagram.png"
	followers := "1000"
	following := "500"
	posts := "200"

	width := 300
	height := 100
	canvas := svg.New(w)
	canvas.Start(width, height)
	defer canvas.End()

	// 背景
	canvas.Rect(0, 0, width, height, "fill:#f0f0f0")

	// アイコン
	canvas.Image(10, 10, 80, 80, iconURL)

	// 統計情報
	canvas.Text(120, 30, fmt.Sprintf("Followers: %s", followers), "font-family:Arial;font-size:14px")
	canvas.Text(120, 55, fmt.Sprintf("Following: %s", following), "font-family:Arial;font-size:14px")
	canvas.Text(120, 80, fmt.Sprintf("Posts: %s", posts), "font-family:Arial;font-size:14px")
}
