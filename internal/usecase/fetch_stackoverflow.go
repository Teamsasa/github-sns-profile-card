package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
)

func FetchStackoverflowData(username string) (*model.PlatformUserInfo, error) {
	for _, c := range username {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("id must be numeric")
		}
	}

	type reputationResult struct {
		Reputation  int
		DisplayName string
	}
	// Buffered results let workers finish even if another request fails first.
	reputationChan := make(chan reputationResult, 1)
	answerCountChan := make(chan int, 1)
	questionCountChan := make(chan int, 1)
	errChan := make(chan error, 3)

	// reputationを取得
	go func(resultChan chan reputationResult) {
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
				Reputation  int    `json:"reputation"`
				DisplayName string `json:"display_name"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respReputation); err != nil {
			errChan <- err
			return
		}
		// ユーザーが削除された場合？にステータスコードは200だが、itemsが空になる
		if len(respReputation.Items) == 0 {
			errChan <- fmt.Errorf("user not found")
			return
		}
		resultChan <- reputationResult{
			Reputation:  respReputation.Items[0].Reputation,
			DisplayName: respReputation.Items[0].DisplayName,
		}
	}(reputationChan)

	// 回答数を取得
	go func(resultChan chan int) {
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
		resultChan <- len(respAnswers.Items)
	}(answerCountChan)

	// 質問数を取得
	go func(resultChan chan int) {
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
		resultChan <- len(respQuestions.Items)
	}(questionCountChan)

	var reputation, answerCount, questionCount int
	var displayName string
	// Completion depends on receiving each result, including valid zero counts.
	for reputationChan != nil || answerCountChan != nil || questionCountChan != nil {
		select {
		case rep := <-reputationChan:
			reputation = rep.Reputation
			displayName = rep.DisplayName
			reputationChan = nil
		case ans := <-answerCountChan:
			answerCount = ans
			answerCountChan = nil
		case ques := <-questionCountChan:
			questionCount = ques
			questionCountChan = nil
		case err := <-errChan:
			return nil, err
		}
	}

	return &model.PlatformUserInfo{
		UserName:      displayName,
		Reputation:    reputation,
		AnswerCount:   answerCount,
		QuestionCount: questionCount,
	}, nil
}

func FormatNumber(n int) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}
