package server

import (
	"fmt"
	"net/http"
	"profile/internal/model"
	"profile/internal/usecase"

	svg "github.com/ajstarks/svgo"
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

	iconURL, exists := model.PlatformIcons[platform]
	if !exists {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform")
		return
	}

	urlBase, exists := model.PlatformURLs[platform]
	if !exists || username == "" {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform or empty username")
		return
	}

	userInfo, err := fetchUserData(platform, username)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data from %s", platform))
		return
	}

	var url string
	if platform == "youtube" {
		url = urlBase + userInfo.CustomURL
	} else {
		url = urlBase + username
	}

	// SVGの生成
	w.Header().Set("Content-Type", "image/svg+xml")

	width := 250
	height := 120
	borderRadius := 20
	strokeWidth := 4
	textColor := model.PlatformFontColors[platform]
	canvas := svg.New(w)
	canvas.Start(width+(2*strokeWidth), height+(2*strokeWidth))
	defer canvas.End()

	// リンクの開始
	canvas.Link(url, "")

	// 外枠と背景（角丸の長方形）を描画
	fmt.Fprintf(canvas.Writer, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" ry="%d" fill="%s" stroke="%s" stroke-width="%d" />`, strokeWidth, strokeWidth, width, height, borderRadius, borderRadius, model.PlatformBgColors[platform], model.PlatformColors[platform], strokeWidth)

	// アイコン
	canvas.Image(20+strokeWidth, 20+strokeWidth, 80, 80, iconURL)

	// 統計情報
	if platform == "stackoverflow" || platform == "note" || platform == "youtube" {
		canvas.Text(130+strokeWidth, 25+strokeWidth, fmt.Sprintf(userInfo.UserName), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 25+strokeWidth, fmt.Sprintf("@%s", username), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}
	if platform == "stackoverflow" {
		canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Reputation: %s", usecase.FormatNumber(userInfo.Reputation)), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "atcoder" {
		canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Ranking: %d", userInfo.Ranking), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Followers: %d", userInfo.FollowersCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Likes: %d", userInfo.LikeCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "stackoverflow" {
		if userInfo.AnswerCount >= 100 {
			canvas.Text(130+strokeWidth, 75+strokeWidth, "Answers: 100+", fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
		} else {
			canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Answers: %d", userInfo.AnswerCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
		}
	} else if platform == "youtube" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Videos: %d", userInfo.TotalVideos), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "atcoder" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Rating: %d", userInfo.Rating), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Following: %d", userInfo.FollowingCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Articles: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "stackoverflow" {
		if userInfo.QuestionCount >= 100 {
			canvas.Text(130+strokeWidth, 100+strokeWidth, "Questions: 100+", fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
		} else {
			canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Questions: %d", userInfo.QuestionCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
		}
	} else if platform == "youtube" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Views: %d", userInfo.TotalViewCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "note" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Notes: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "atcoder" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Matches: %d", userInfo.RatedMatches), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Posts: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	// リンクの終了
	canvas.LinkEnd()
}

// 各プラットフォームからデータを取得する関数
func fetchUserData(platform, username string) (*model.PlatformUserInfo, error) {
	switch platform {
	case "qiita":
		return usecase.FetchQiitaData(username)
	case "twitter":
		return usecase.FetchTwitterData(username)
	case "zenn":
		return usecase.FetchZennData(username)
	case "linkedin":
		return usecase.FetchLinkedinData(username)
	case "stackoverflow":
		return usecase.FetchStackoverflowData(username)
	case "atcoder":
		return usecase.FetchAtCoderData(username)
	case "note":
		return usecase.FetchNoteData(username)
	case "youtube":
		return usecase.FetchYoutubeData(username)
	case "instagram":
		return usecase.FetchInstagramData(username)
	}
	return nil, fmt.Errorf("platform not supported")
}
