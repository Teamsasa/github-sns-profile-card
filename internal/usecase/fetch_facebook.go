package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"profile/internal/model"
)

func FetchFacebookData(userID string) (*model.PlatformUserInfo, error) {
	accessToken := os.Getenv("FACEBOOK_ACCESS_TOKEN")

	url := fmt.Sprintf("https://graph.facebook.com/v17.0/%s?fields=followers_count,friends_count,posts.limit(0).summary(true)&access_token=%s", userID, accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    int    `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("Facebook API error: %s (type: %s, code: %d)", errorResp.Error.Message, errorResp.Error.Type, errorResp.Error.Code)
		}// エラー詳細を取得するために追加　怒られてるけど
		return nil, fmt.Errorf("API error: status code %d, response: %s", resp.StatusCode, string(body))
	}

	var user struct {
		FollowersCount int `json:"followers_count"`
		FriendsCount   int `json:"friends_count"`
		Posts struct {
			Summary struct {
				TotalCount int `json:"total_count"`
			} `json:"summary"`
		} `json:"posts"`
	}

	if err := json.Unmarshal(body, &user); err !=nil {
		return nil, fmt.Errorf("JSON decode error: %w", err)
	}

	return &model.PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FriendsCount,
		ArticlesCount:  user.Posts.Summary.TotalCount,
	}, nil
}
