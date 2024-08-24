package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
)

// Qiitaのユーザーデータを取得する関数
func FetchQiitaData(username string) (*model.PlatformUserInfo, error) {
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

	return &model.PlatformUserInfo{
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FolloweesCount,
		ArticlesCount:  user.ArticlesCount,
	}, nil
}
