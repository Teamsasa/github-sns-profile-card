package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"profile/internal/model"
	"sync"
)

func FetchStackoverflowData(username string) (*model.PlatformUserInfo, error) {
	for _, c := range username {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("id must be numeric")
		}
	}

	var wg sync.WaitGroup
	reputationChan := make(chan struct {
		Reputation  int
		DisplayName string
	})
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
		reputationChan <- struct {
			Reputation  int
			DisplayName string
		}{
			Reputation:  respReputation.Items[0].Reputation,
			DisplayName: respReputation.Items[0].DisplayName,
		}
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
	var displayName string
	for {
		select {
		case rep, ok := <-reputationChan:
			if ok {
				reputation = rep.Reputation
				displayName = rep.DisplayName
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
