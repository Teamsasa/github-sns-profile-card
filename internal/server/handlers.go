package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

var platformIcons = map[string]string{
	"zenn":          "/assets/zenn.png",
	"qiita":         "/assets/qiita.png",
	"twitter":       "/assets/twitter.png",
	"linkedin":      "/assets/linkedin.png",
	"stackoverflow": "/assets/stackoverflow.png",
	"atcoder":       "/assets/atcoder.png",
}

type PlatformUserInfo struct {
	FollowersCount int
	FollowingCount int
	ArticlesCount  int
	Rating         int // AtCoder用のフィールド
}

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

	iconURL, exists := platformIcons[platform]
	if !exists {
		handleError(w, nil, http.StatusBadRequest, "Unknown platform")
		return
	}

	userInfo, err := fetchUserData(platform, username)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data from %s", platform))
		return
	}

	// SVGの生成
	w.Header().Set("Content-Type", "image/svg+xml")

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
	canvas.Text(120, 30, fmt.Sprintf("Followers: %d", userInfo.FollowersCount), "font-family:Arial;font-size:14px")
	canvas.Text(120, 55, fmt.Sprintf("Following: %d", userInfo.FollowingCount), "font-family:Arial;font-size:14px")
	canvas.Text(120, 80, fmt.Sprintf("Posts: %d", userInfo.ArticlesCount), "font-family:Arial;font-size:14px")

	// AtCoderの場合はRatingも表示
	if platform == "atcoder" {
		canvas.Text(120, 105, fmt.Sprintf("Rating: %d", userInfo.Rating), "font-family:Arial;font-size:14px")
	}
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

// Qiitaのユーザーデータを取得する関数
func fetchQiitaData(username string) (*PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://qiita.com/api/v2/users/%s", username))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		FollowersCount int `json:"followers_count"`
		FolloweesCount int `json:"followees_count"`
		ArticlesCount  int `json:"items_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FolloweesCount,
		ArticlesCount:  user.ArticlesCount,
	}, nil
}

func fetchTwitterData(username string) (*PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", username))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		FollowersCount int `json:"followers_count"`
		FolloweesCount int `json:"following_count"` // Twitter APIでは"following_count"と呼ばれる場合が多い
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FolloweesCount,
		ArticlesCount:  0, // Twitterは投稿数の取得がAPIでサポートされていないため、0を返します
	}, nil
}

func fetchZennData(username string) (*PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://zenn.dev/api/users/%s", username))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		FollowersCount int `json:"followers_count"`
		ArticlesCount  int `json:"articles_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: 0, // Zenn APIにはフォロー中のユーザー数がないため、0を返します
		ArticlesCount:  user.ArticlesCount,
	}, nil
}

func fetchLinkedinData(username string) (*PlatformUserInfo, error) {
	//課金が必要なため、実装は省略
	_ = username
	return nil, fmt.Errorf("not implemented")
}

func fetchStackoverflowData(username string) (*PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.stackexchange.com/2.3/users/%s?site=stackoverflow", username))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Items []struct {
			Reputation   int `json:"reputation"`
			AnswerCount  int `json:"answer_count"`
			QuestionCount int `json:"question_count"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := response.Items[0]

	return &PlatformUserInfo{
		FollowersCount: user.Reputation,   // StackOverflowではReputationをFollowersCountとして代用
		FollowingCount: 0,                 // StackOverflow APIにはフォロー中のユーザー数がないため、0を返します
		ArticlesCount:  user.AnswerCount + user.QuestionCount, // 回答数と質問数の合計を投稿数として扱います
	}, nil
}

// AtCoderのユーザーデータを取得する関数
func fetchAtCoderData(username string) (*PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://atcoder.jp/users/%s", username)) // 仮のURL
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		Rating int `json:"rating"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &PlatformUserInfo{
		Rating: user.Rating,
	}, nil
}
