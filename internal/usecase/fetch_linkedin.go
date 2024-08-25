package usecase

import (
	"fmt"
	"net/http"
	"profile/internal/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchLinkedinData(username string) (*model.PlatformUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.linkedin.com/in/%s?original_referer=https://www.google.com/", username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data")
	}

	return parseLinkedinHTML(resp)
}

func parseLinkedinHTML(resp *http.Response) (*model.PlatformUserInfo, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	selection := doc.Find("meta[name=description]")
	if selection.Length() == 0 {
		return nil, fmt.Errorf("description attribute not found")
	}
	description, exists := selection.Attr("content")
	if !exists {
		return nil, fmt.Errorf("content attribute not found")
	}

	splitted := strings.Split(description, " · ")
	if len(splitted) < 3 {
		return nil, fmt.Errorf("failed to split description")
	}
	experience := strings.TrimPrefix(splitted[1], "Experience: ")
	education := strings.TrimPrefix(splitted[2], "Education: ")
	location := strings.TrimPrefix(splitted[3], "Location: ")

	selection = doc.Find("meta[property='profile:first_name']")
	firstName, exists := selection.Attr("content")
	if !exists {
		return nil, fmt.Errorf("first_name attribute not found")
	}
	selection = doc.Find("meta[property='profile:last_name']")
	lastName, exists := selection.Attr("content")
	if !exists {
		return nil, fmt.Errorf("last_name attribute not found")
	}

	return &model.PlatformUserInfo{
		UserName:   firstName + " " + lastName,
		Experience: experience,
		Education:  education,
		Location:   location,
	}, nil
}
