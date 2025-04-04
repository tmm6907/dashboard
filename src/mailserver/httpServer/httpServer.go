package httpServer

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the RSS feed details
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

func MakeRSSFeedBody(feedName, link string, size int) *RSS {
	return &RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       feedName,
			Description: "Email Feed",
			Link:        link,
			Items:       make([]Item, size),
		},
	}
}

// Item represents an individual RSS item (email entry)
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Mail struct {
	ID        uint   `db:"id"`
	Sender    string `db:"sender"`
	Recipient string `db:"recipient"`
	Subject   string `db:"subject"`
	Body      string `db:"body"`
	CreatedAt string `db:"created_at"`
}

func generateXML(feedName string, recipient string, link string) ([]byte, error) {
	var mail []Mail
	db, err := sqlx.Open("sqlite3", "mail.db")
	if err != nil {
		return nil, err
	}
	if err = db.Select(&mail, "SELECT * FROM mail WHERE recipient = ?;", recipient); err != nil {
		return nil, err
	}
	rssFeed := MakeRSSFeedBody(feedName, link, len(mail))
	for i, item := range mail {
		item := Item{
			Title:       item.Subject,
			Description: item.Body,
			PubDate:     item.CreatedAt,
		}
		rssFeed.Channel.Items[i] = item
	}

	xmlData, err := xml.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Info("Generated RSS feed for ", link)
	return xmlData, nil
}

func FeedHandler(c *fiber.Ctx) error {
	var reqBody struct {
		FeedName  string `json:"feedName"`
		Recipient string `json:"recipient"`
		Link      string `json:"link"`
	}
	if err := json.Unmarshal(c.Body(), &reqBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	responseBytes, err := generateXML(reqBody.FeedName, reqBody.Recipient, reqBody.Link)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	c.Set("Content-Type", "application/xml")
	return c.Send(responseBytes)
}
