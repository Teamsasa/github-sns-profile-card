package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"profile/internal/model"
)

func FetchInstagramData(userID string) (*model.PlatformUserInfo, error) {
	accessToken := os.Getenv("FACEBOOK_ACCESS_TOKEN")
	// Instagram Graph APIのエンドポイントを設定
	url := fmt.Sprintf("https://graph.instagram.com/%s?fields=followers_count,follows_count,media_count&access_token=%s", userID, accessToken)

	// HTTP GETリクエストを送信
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードがOKでない場合の処理
if resp.StatusCode != http.StatusOK {
    var errorResp struct {
        Error struct {
            Message string `json:"message"`
            Type    string `json:"type"`
            Code    int    `json:"code"`
            ErrorUserTitle string `json:"error_user_title"` // エラー詳細を取得するために追加
            ErrorUserMsg string `json:"error_user_msg"`   // エラー詳細を取得するために追加
        } `json:"error"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
        return nil, fmt.Errorf("Instagram API error: %s (type: %s, code: %d, user_title: %s, user_msg: %s)", errorResp.Error.Message, errorResp.Error.Type, errorResp.Error.Code, errorResp.Error.ErrorUserTitle, errorResp.Error.ErrorUserMsg)
    }// エラー詳細を取得するために追加　怒られてるけど
    return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
}


	// レスポンスデータを解析
	var user struct {
		FollowersCount int `json:"followers_count"`
		FollowingCount int `json:"follows_count"`
		PostCount      int `json:"media_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}

	// 結果を返す
	return &model.PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		ArticlesCount:  user.PostCount,
	}, nil
}
