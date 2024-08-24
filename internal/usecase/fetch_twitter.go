package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
"profile/internal/model")

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
