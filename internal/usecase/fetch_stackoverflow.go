package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
)

func FetchStackoverflowData(username string) (*model.PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.stackexchange.com/2.3/users/%s?site=stackoverflow", username))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Items []struct {
			Reputation    int `json:"reputation"`
			AnswerCount   int `json:"answer_count"`
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

	return &model.PlatformUserInfo{
		FollowersCount: user.Reputation,                       // StackOverflowではReputationをFollowersCountとして代用
		FollowingCount: 0,                                     // StackOverflow APIにはフォロー中のユーザー数がないため、0を返します
		ArticlesCount:  user.AnswerCount + user.QuestionCount, // 回答数と質問数の合計を投稿数として扱います
	}, nil
}
