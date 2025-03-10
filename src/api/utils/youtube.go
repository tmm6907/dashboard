package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

func ExtractChannelID(youtubeURL string) (string, error) {
	resp, err := http.Get(youtubeURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch page: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read page body: %v", err)
	}

	// Regex to find the channel ID in og:url meta tag
	re := regexp.MustCompile(`https://www\.youtube\.com/channel/([A-Za-z0-9_-]+)`)
	match := re.FindStringSubmatch(string(body))

	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("channel ID not found")
}

func IsYoutubeChannelURL(url string) bool {
	return strings.Contains(url, "www.youtube.com")
}

func getFeedFromChannelID(channelID string) string {
	return "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID
}

func GetYouTubeRSS(channelURL string) (string, error) {
	handleRegex := regexp.MustCompile(`https?://(?:www\.)?youtube\.com/@([a-zA-Z0-9_-]+)`)
	customRegex := regexp.MustCompile(`https?://(www\.)?youtube\.com/c/[a-zA-Z0-9_-]+`)
	channelIDRegex := regexp.MustCompile(`https?://(www\.)?youtube\.com/channel/[A-Za-z0-9_-]+`)
	userRegex := regexp.MustCompile(`https?://(www\.)?youtube\.com/user/[a-zA-Z0-9_-]+`)
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "", errors.New("missing api key")
	}
	var channelID string
	switch {
	case handleRegex.MatchString(channelURL):
		//https://www.googleapis.com/youtube/v3/channels?part=id&forUsername=@LegalEagle&key=
		log.Debug(handleRegex.FindStringSubmatch(channelURL))
		identifier := handleRegex.FindStringSubmatch(channelURL)[1]
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=id&forHandle=@%s&key=%s", identifier, apiKey)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data := make(map[string]any)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", err
		}

		items, ok := data["items"].([]interface{})
		if !ok {
			return "", fmt.Errorf("Invalid response format: items is not a list, got %T", data["items"])
		}

		if len(items) != 1 {
			return "", fmt.Errorf("Unknown exception: %v", data)
		}
		item, ok := items[0].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("Invalid response format: item is not a map, got %T", items[0])
		}

		id, ok := item["id"].(string)
		if !ok {
			return "", fmt.Errorf("Invalid response format: id is not a string, got %T", item["id"])
		}

		channelID = id

	case customRegex.MatchString(channelURL):
		identifier, err := ExtractChannelID(channelURL)
		if err != nil {
			return "", err
		}
		channelID = identifier
		break
	case channelIDRegex.MatchString(channelURL):
		channelID = channelIDRegex.FindStringSubmatch(channelURL)[2]
		break
	case userRegex.MatchString(channelURL):
		identifier := userRegex.FindStringSubmatch(channelURL)[2]
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=id&forUsername=%s&key=%s", identifier, apiKey)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data := make(map[string]any)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return "", err
		}

		items := data["items"].([]map[string]any)

		if len(items) != 1 {
			return "", fmt.Errorf("Unknown exception: %v", data)
		}
		channelID = items[0]["id"].(string)
		break

	default:
		return "", fmt.Errorf("Unrecognized channelURL: %s", channelURL)
	}

	if channelID == "" {
		return "", fmt.Errorf("Unrecognized channelURL: %s", channelURL)
	}
	return getFeedFromChannelID(channelID), nil
}
