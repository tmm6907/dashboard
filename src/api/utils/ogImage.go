package utils

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func GetOGImage(articleURL string) (string, error) {
	client := &http.Client{
		// Follow redirects manually
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent automatic redirect so we can capture the real location
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(articleURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		loc, err := resp.Location()
		if err != nil {
			return "", fmt.Errorf("redirect location error: %v", err)
		}
		articleURL = loc.String()
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch URL, status: %d", resp.StatusCode)
	}

	// Step 2: Now fetch the real article page
	resp, err = http.Get(articleURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the HTML
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			// End of document
			return "", errors.New("og:image not found")
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				// Check attributes for property="og:image"
				var content string
				for _, attr := range token.Attr {
					if attr.Key == "property" && attr.Val == "og:image" {
						// Once we find property="og:image", get the content attribute
						for _, attr := range token.Attr {
							if attr.Key == "content" {
								content = attr.Val
								return content, nil
							}
						}
					}
				}
			}
		}
	}
}
