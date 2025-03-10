package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
)

// GetLoginURL returns the OAuth login URL
func GetLoginURL(config *oauth2.Config, state string) string {
	return config.AuthCodeURL(state)
}

// GetUserInfo exchanges code for user info
func GetUserInfo(config *oauth2.Config, code string) (map[string]any, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var userInfo map[string]any
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
