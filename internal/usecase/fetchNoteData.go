package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
)

func FetchNoteData(username string) (*model.PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://note.com/api/v2/creators/%s", username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data")
	}

	var user struct {
		Data struct {
			FollowingCount int    `json:"followingCount"`
			FollowersCount int    `json:"followerCount"`
			AritclesCount  int    `json:"noteCount"`
			Nickname       string `json:"nickname"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &model.PlatformUserInfo{
		FollowersCount: user.Data.FollowersCount,
		FollowingCount: user.Data.FollowingCount,
		ArticlesCount:  user.Data.AritclesCount,
		UserName:       user.Data.Nickname,
	}, nil
}
