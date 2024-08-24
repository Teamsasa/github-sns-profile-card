package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"profile/internal/model"
)

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
