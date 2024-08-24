package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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

var platformURLs = map[string]string{
	"zenn":          "https://zenn.dev/",
	"qiita":         "https://qiita.com/",
	"twitter":       "https://x.com/",
	"linkedin":      "https://linkedin.com/in/",
	"stackoverflow": "https://stackoverflow.com/users/",
	"atcoder":       "https://atcoder.jp/usres/",
}

var platformColors = map[string]string{
	"zenn":          "#3EA8FF",
	"qiita":         "#55C500",
	"twitter":       "#FFFFFF",
	"linkedin":      "#0A66C2",
	"stackoverflow": "#F48024",
	"atcoder":       "#000000",
}

var platformBgColors = map[string]string{
	"zenn":          "#F1F5F9",
	"qiita":         "#F5F6F6",
	"twitter":       "#000000",
	"linkedin":      "##F4F2EE",
	"stackoverflow": "#FFFFFB",
	"atcoder":       "#EBEBEB",
}

var platformFontColors = map[string]string{
	"zenn":          "#000000",
	"qiita":         "#000000",
	"twitter":       "#FFFFFF",
	"linkedin":      "#000000",
	"stackoverflow": "#000000",
	"atcoder":       "#000000",
}

type PlatformUserInfo struct {
	FollowersCount int
	FollowingCount int
	ArticlesCount  int
	LikeCount      int // Zenn用のフィールド
	Reputation     int // StackOverflow用のフィールド
	AnswerCount    int // StackOverflow用のフィールド
	QuestionCount  int // StackOverflow用のフィールド
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

	urlBase, exists := platformURLs[platform]
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
	textColor := platformFontColors[platform]
	canvas := svg.New(w)
	canvas.Start(width+(2*strokeWidth), height+(2*strokeWidth))
	defer canvas.End()

	// リンクの開始
	canvas.Link(url, "")

	// 外枠を描画
	canvas.Rect(strokeWidth, strokeWidth, width, height,
		fmt.Sprintf("fill:none;rx:%d;ry:%d;stroke:%s;stroke-width:%d", borderRadius, borderRadius, platformColors[platform], strokeWidth))

	// 背景（角丸の長方形）
	canvas.Rect(strokeWidth, strokeWidth, width, height,
		fmt.Sprintf("fill:%s;rx:%d;ry:%d", platformBgColors[platform], borderRadius, borderRadius))

	// アイコン
	canvas.Image(20+strokeWidth, 20+strokeWidth, 80, 80, iconURL)

	// 統計情報
	canvas.Text(130+strokeWidth, 25+strokeWidth, fmt.Sprintf("@%s", username), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	if platform == "stackoverflow" {
		canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Reputation: %d", userInfo.Reputation), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 50+strokeWidth, fmt.Sprintf("Followers: %d", userInfo.FollowersCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Likes: %d", userInfo.LikeCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "stackoverflow" {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Answers: %d", userInfo.AnswerCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else {
		canvas.Text(130+strokeWidth, 75+strokeWidth, fmt.Sprintf("Following: %d", userInfo.FollowingCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	}

	if platform == "zenn" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Articles: %d", userInfo.ArticlesCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
	} else if platform == "stackoverflow" {
		canvas.Text(130+strokeWidth, 100+strokeWidth, fmt.Sprintf("Questions: %d", userInfo.QuestionCount), fmt.Sprintf("font-family:Arial;font-size:14px;fill:%s", textColor))
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
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user not found")
	}
	defer resp.Body.Close()

	var user struct {
		User struct {
			FollowersCount int `json:"follower_count"`
			LikeCount      int `json:"total_liked_count"`
			ArticlesCount  int `json:"articles_count"`
		} `json:"user"`
	}

	fmt.Println(resp.Body)
	fmt.Println(user)

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &PlatformUserInfo{
		FollowersCount: user.User.FollowersCount,
		LikeCount:      user.User.LikeCount,
		ArticlesCount:  user.User.ArticlesCount,
	}, nil
}

func fetchLinkedinData(username string) (*PlatformUserInfo, error) {
	//課金が必要なため、実装は省略
	_ = username
	return nil, fmt.Errorf("not implemented")
}

func fetchStackoverflowData(username string) (*PlatformUserInfo, error) {
	for _, c := range username {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("id must be numeric")
		}
	}

	var wg sync.WaitGroup
	reputationChan := make(chan int)
	answerCountChan := make(chan int)
	questionCountChan := make(chan int)
	errChan := make(chan error, 3)

	// reputationを取得
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("https://api.stackexchange.com/2.3/users/%s?site=stackoverflow", username))
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("fetch failed")
			return
		}
		var respReputation struct {
			Items []struct {
				Reputation int `json:"reputation"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respReputation); err != nil {
			errChan <- err
			return
		}
		reputationChan <- respReputation.Items[0].Reputation
	}()

	// 回答数を取得
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("https://api.stackexchange.com/2.3/users/%s/answers?pagesize=100&site=stackoverflow", username))
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("fetch failed")
			return
		}
		var respAnswers struct {
			Items []struct {
				Content []interface{} `json:"content"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respAnswers); err != nil {
			errChan <- err
			return
		}
		answerCountChan <- len(respAnswers.Items)
	}()

	// 質問数を取得
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get(fmt.Sprintf("https://api.stackexchange.com/2.3/users/%s/questions?pagesize=100&site=stackoverflow", username))
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("fetch failed")
			return
		}
		var respQuestions struct {
			Items []struct {
				Content []interface{} `json:"content"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respQuestions); err != nil {
			errChan <- err
			return
		}
		questionCountChan <- len(respQuestions.Items)
	}()

	go func() {
		wg.Wait()
		close(reputationChan)
		close(answerCountChan)
		close(questionCountChan)
		close(errChan)
	}()

	var reputation, answerCount, questionCount int
	for {
		select {
		case rep, ok := <-reputationChan:
			if ok {
				reputation = rep
			}
		case ans, ok := <-answerCountChan:
			if ok {
				answerCount = ans
			}
		case ques, ok := <-questionCountChan:
			if ok {
				questionCount = ques
			}
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		}
		if reputation != 0 && answerCount != 0 && questionCount != 0 {
			break
		}
	}

	return &PlatformUserInfo{
		Reputation:    reputation,
		AnswerCount:   answerCount,
		QuestionCount: questionCount,
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
