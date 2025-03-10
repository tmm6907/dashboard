package worker

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mmcdole/gofeed"
	"github.com/tmm6907/dashboard/models"
)

func (h *Handler) FecthRSSFeed(feed models.Feed) error {
	rssParser := gofeed.NewParser()
	rssFeed, err := rssParser.ParseURL(feed.Link)
	if err != nil {
		return err
	}
	db := h.GetDB()
	for _, item := range rssFeed.Items {
		image := ""
		if item.Image != nil {
			image = item.Image.URL
		}
		if _, err = db.Exec("INSERT OR IGNORE INTO feed_items(feed_id, title, link, description, image, guid, pub_date) VALUES (?, ?, ?, ?, ?, ?, ?);", feed.FeedID, item.Title, item.Link, item.Description, image, item.GUID, item.Published); err != nil {
			return err
		}

	}

	return nil
}

func (h *Handler) FecthRSSFeeds() {
	var feeds []models.Feed
	db := h.GetDB()

	err := db.Select(&feeds, "SELECT * FROM feeds;")
	if err != nil {
		log.Error(err)
	}
	workers := 10
	feedChan := make(chan models.Feed, len(feeds))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range feedChan {
				log.Info("Fetching RSS Feed for url: ", f.Link, " ...")
				if err := h.FecthRSSFeed(f); err != nil {
					log.Error(err)
				}
			}
		}()
	}
	for _, feed := range feeds {
		feedChan <- feed
	}
	close(feedChan)
	wg.Wait()
}

func (h *Handler) StartRSSFetcher(interval *time.Duration) {
	if interval == nil {
		defaultDuration := 1 * time.Minute
		interval = &defaultDuration
	}
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info("Fetching RSS Feeds...")
			h.FecthRSSFeeds()
		}
	}
}
