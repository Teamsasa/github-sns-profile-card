package server

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.SVGHandler)

	return mux
}

func (s *Server) SVGHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	iconURL := "https://static.cdninstagram.com/rsrc.php/v3/yI/r/VsNE-OHk_8a.png"
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
