package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
)

func FetchZennData(username string) (*model.PlatformUserInfo, error) {
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

	return &model.PlatformUserInfo{
		FollowersCount: user.User.FollowersCount,
		LikeCount:      user.User.LikeCount,
		ArticlesCount:  user.User.ArticlesCount,
	}, nil
}
