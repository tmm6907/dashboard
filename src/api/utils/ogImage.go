package utils

import (
	"errors"
	"net/http"

	"golang.org/x/net/html"
)

func GetOGImage(url string) (string, error) {
	resp, err := http.Get(url)
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
