package usecase

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"profile/internal/model"
)

// AtCoderのユーザーデータを取得する関数
func FetchAtCoderData(username string) (*model.PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://atcoder.jp/users/%s", username)) // 仮のURL
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data")
	}

	user, err := parseAtCoderHTML(resp)

	if err != nil {
		return nil, err
	}

	return &model.PlatformUserInfo{
		Ranking:      user.Ranking,
		Rating:       user.Rating,
		RatedMatches: user.RatedMatches,
	}, nil
}

func parseAtCoderHTML(resp *http.Response) (*model.PlatformUserInfo, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	rankingStr := doc.Find("div.row div.col-md-9 table tr").Eq(0).Find("td").Text()
	ratingStr := doc.Find("div.row div.col-md-9 table tr").Eq(1).Find("span").Text()
	ratedMatchesStr := doc.Find("div.row div.col-md-9 table tr").Eq(3).Find("td").Text()

	if rankingStr == "" || ratingStr == "" || ratedMatchesStr == "" {
		return nil, fmt.Errorf("failed to parse data")
	}

	ranking, err := strconv.Atoi(rankingStr[:len(rankingStr)-2])
	if err != nil {
		return nil, err
	}
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		return nil, err
	}
	ratedMatches, err := strconv.Atoi(ratedMatchesStr)
	if err != nil {
		return nil, err
	}

	return &model.PlatformUserInfo{
		Ranking:      ranking,
		Rating:       rating,
		RatedMatches: ratedMatches,
	}, nil
}
