package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	cardWidth := r.URL.Query().Get("width")
	if platform == "" || username == "" {
		handleError(w, nil, http.StatusBadRequest, "platform and userid are required")
		return
	}

	iconURL, exists := model.PlatformIcons[platform]
	if !exists {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform")
		return
	}
	// アイコンをローカルファイルから取得してBase64エンコード
	f, err := os.Open(iconURL)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to open icon file")
		return
	}
	defer f.Close()
	fileStat, err := f.Stat()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to get icon file info")
		return
	}
	iconData := make([]byte, fileStat.Size())
	f.Read(iconData)
	iconDataBase64 := base64.StdEncoding.EncodeToString(iconData)

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

	var width int
	if len(cardWidth) > 0 {
		width, err = strconv.Atoi(cardWidth)
		if err != nil {
			handleError(w, err, http.StatusBadRequest, "Invalid width")
			return
		}
	} else {
	width = 260
	}
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
	canvas.Image(15+strokeWidth, 20+strokeWidth, 80, 80, "data:image/png;base64,"+iconDataBase64)

	fontStyle := fmt.Sprintf("font-family:opensans;font-size:15px;fill:%s;font-weight:bold;", textColor)

	// 統計情報
	if platform == "stackoverflow" || platform == "note" || platform == "youtube" {
		canvas.Text(110+strokeWidth, 25+strokeWidth, fmt.Sprintf(userInfo.UserName), fontStyle)
	} else {
		canvas.Text(110+strokeWidth, 25+strokeWidth, fmt.Sprintf("@%s", username), fontStyle)
	}
	if platform == "stackoverflow" {
		canvas.Text(110+strokeWidth, 50+strokeWidth, fmt.Sprintf("Reputation: %s", usecase.FormatNumber(userInfo.Reputation)), fontStyle)
	} else if platform == "atcoder" {
		canvas.Text(110+strokeWidth, 50+strokeWidth, fmt.Sprintf("Ranking: %s", usecase.FormatNumber(userInfo.Ranking)), fontStyle)
	} else {
		canvas.Text(110+strokeWidth, 50+strokeWidth, fmt.Sprintf("Followers: %s", usecase.FormatNumber(userInfo.FollowersCount)), fontStyle)
	}

	if platform == "zenn" {
		canvas.Text(110+strokeWidth, 75+strokeWidth, fmt.Sprintf("Likes: %s", usecase.FormatNumber(userInfo.LikeCount)), fontStyle)
	} else if platform == "stackoverflow" {
		if userInfo.AnswerCount >= 100 {
			canvas.Text(110+strokeWidth, 75+strokeWidth, "Answers: 100+", fontStyle)
		} else {
			canvas.Text(110+strokeWidth, 75+strokeWidth, fmt.Sprintf("Answers: %s", usecase.FormatNumber(userInfo.AnswerCount)), fontStyle)
		}
	} else if platform == "youtube" {
		canvas.Text(110+strokeWidth, 75+strokeWidth, fmt.Sprintf("Videos: %s", usecase.FormatNumber(userInfo.TotalVideos)), fontStyle)
	} else if platform == "atcoder" {
		canvas.Text(110+strokeWidth, 75+strokeWidth, fmt.Sprintf("Rating: %s", usecase.FormatNumber(userInfo.Rating)), fontStyle)
	} else {
		canvas.Text(110+strokeWidth, 75+strokeWidth, fmt.Sprintf("Following: %s", usecase.FormatNumber(userInfo.FollowingCount)), fontStyle)
	}

	if platform == "zenn" {
		canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Articles: %s", usecase.FormatNumber(userInfo.ArticlesCount)), fontStyle)
	} else if platform == "stackoverflow" {
		if userInfo.QuestionCount >= 100 {
			canvas.Text(110+strokeWidth, 100+strokeWidth, "Questions: 100+", fontStyle)
		} else {
			canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Questions: %s", usecase.FormatNumber(userInfo.QuestionCount)), fontStyle)
		}
	} else if platform == "youtube" {
		canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Views: %s", usecase.FormatNumber(userInfo.TotalViewCount)), fontStyle)
	} else if platform == "note" {
		canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Notes: %s", usecase.FormatNumber(userInfo.ArticlesCount)), fontStyle)
	} else if platform == "atcoder" {
		canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Matches: %s", usecase.FormatNumber(userInfo.RatedMatches)), fontStyle)
	} else {
		canvas.Text(110+strokeWidth, 100+strokeWidth, fmt.Sprintf("Posts: %s", usecase.FormatNumber(userInfo.ArticlesCount)), fontStyle)
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
