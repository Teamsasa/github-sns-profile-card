package server

import (
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"profile/internal/usecase"
	"profile/internal/model"
)


// 汎用エラーハンドリング関数
func handleError(w http.ResponseWriter, err error, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
	fmt.Println("Error:", err)
}

// SVGHandlerは、指定されたプラットフォームのユーザーデータを取得し、SVG画像を生成するハンドラー
func (s *Server) SVGHandler(w http.ResponseWriter, r *http.Request) {
	platform := r.URL.Query().Get("platform")
	username := r.URL.Query().Get("userid")
	if platform == "" || username == "" {
		handleError(w, nil, http.StatusBadRequest, "platform and userid are required")
		return
	}

	iconURL, exists := PlatformIcons[platform]
	if !exists {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform")
		return
	}

	urlBase, exists := PlatformURLs[platform]
	if !exists || username == "" {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform or empty username")
		return
	}

	url := urlBase + username

	userInfo, err := fetchUserData(platform, username)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data from %s", platform))
		return
	}

	// SVGの生成
	w.Header().Set("Content-Type", "image/svg+xml")

	width := 250
	height := 120
	borderRadius := 20
	strokeWidth := 4
	textColor := PlatformFontColors[platform]
	canvas := svg.New(w)
	canvas.Start(width+(2*strokeWidth), height+(2*strokeWidth))
	defer canvas.End()

	// リンクの開始
	canvas.Link(url, "")

	// 外枠を描画
	canvas.Rect(strokeWidth, strokeWidth, width, height,
		fmt.Sprintf("fill:none;rx:%d;ry:%d;stroke:%s;stroke-width:%d", borderRadius, borderRadius, PlatformColors[platform], strokeWidth))

	// 背景（角丸の長方形）
	canvas.Rect(strokeWidth, strokeWidth, width, height,
		fmt.Sprintf("fill:%s;rx:%d;ry:%d", platformBgColors[platform], borderRadius, borderRadius))

	// アイコン
	canvas.Image(20+strokeWidth, 20+strokeWidth, 80, 80, iconURL)

	// 統計情報
	canvas.Text(130+strokeWidth, 25+strokeWidth, fmt.Sprintf("@%s", username), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Followers: %d", userInfo.FollowersCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Likes: %d", userInfo.LikeCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Following: %d", userInfo.FollowingCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Articles: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Posts: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	// AtCoderの場合はRatingも表示
	if platform == "atcoder" {
		canvas.Text(120+strokeWidth, 130+strokeWidth, fmt.Sprintf("Rating: %d", userInfo.Rating), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	// リンクの終了
	canvas.LinkEnd()
}

// 各プラットフォームからデータを取得する関数
func fetchUserData(platform, username string) (*PlatformUserInfo, error) {
	switch platform {
	case "qiita":
		return fetchQiitaData(username)
	case "twitter":
		return fetchTwitterData(username)
	case "zenn":
		return fetchZennData(username)
	case "linkedin":
		return fetchLinkedinData(username)
	case "stackoverflow":
		return fetchStackoverflowData(username)
	case "atcoder":
		return fetchAtCoderData(username)
	}
	return nil, fmt.Errorf("platform not supported")
}
