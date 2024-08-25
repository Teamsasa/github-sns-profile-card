package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"profile/internal/model"
)

func FetchYoutubeData(username string) (*model.PlatformUserInfo, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is not set in .env")
	}

	channelName := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet&id=%s&key=%s", username, apiKey)
	snippetResponse, err := http.Get(channelName)
	if err != nil {
		return nil, err
	}
	defer snippetResponse.Body.Close()

	if snippetResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", snippetResponse.Status)
	}
	channelURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=statistics&id=%s&key=%s", username, apiKey)
	statisticsResponse, err := http.Get(channelURL)
	if err != nil {
		return nil, err
	}
	defer statisticsResponse.Body.Close()

	if statisticsResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", statisticsResponse.Status)
	}

	var result struct {
		Items []struct {
			Snippet struct {
				Username  string `json:"title"`
				CustomUrl string `json:"customUrl"`
			} `json:"snippet"`
			Statistics struct {
				SubscriberCount string `json:"subscriberCount"`
				ViewCount       string `json:"viewCount"`
				VideoCount      string `json:"videoCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(snippetResponse.Body).Decode(&result); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(statisticsResponse.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no channel found for username: %s", username)
	}

	stats := result.Items[0].Statistics

	customURL := result.Items[0].Snippet.CustomUrl
	username = result.Items[0].Snippet.Username
	subscriberCount, _ := strconv.Atoi(stats.SubscriberCount)
	viewCount, _ := strconv.Atoi(stats.ViewCount)
	videoCount, _ := strconv.Atoi(stats.VideoCount)

	platformUserInfo := &model.PlatformUserInfo{
		CustomURL:      customURL,
		UserName:       username,
		FollowersCount: subscriberCount,
		TotalVideos:    videoCount,
		TotalViewCount: viewCount,
	}

	return platformUserInfo, nil
}
