package worker

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mmcdole/gofeed"
	"github.com/tmm6907/dashboard/models"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) FetchRSSFeed(feed models.Feed) error {
	rssParser := gofeed.NewParser()
	rssFeed, err := rssParser.ParseURL(feed.Link)
	if err != nil {
		return err
	}

	db := h.GetDB()
	feedImage := ""
	feedAlt := ""
	// category, itunes category
	feedCategories := ""

	if rssFeed.Image != nil && feed.Image == "" {
		feedImage = rssFeed.Image.URL
		feedAlt = rssFeed.Image.Title
		if _, err := db.Exec("UPDATE feeds SET image = ?, alt_text = ? WHERE feed_id = ?;", feedImage, feedAlt, feed.FeedID); err != nil {
			return err
		}
	}
	if len(rssFeed.Categories) > 0 && feed.Categories == "" {
		feedCategories = strings.Join(rssFeed.Categories, ", ")
		if _, err := db.Exec("UPDATE feeds SET categories = ? WHERE feed_id = ?;", feedCategories, feed.FeedID); err != nil {
			return err
		}
	}
	mediaType := ""

	if feed.MediaType != "" {
		mediaType = feed.MediaType
	} else if rssFeed.FeedType != "" {
		mediaType = rssFeed.FeedType
	}

	for _, item := range rssFeed.Items {
		image := feedImage
		alt := feedAlt
		categories := ""
		media := ""
		if len(item.Enclosures) > 0 {
			media = item.Enclosures[0].Type
		}
		if item.Image != nil {
			image = item.Image.URL
			alt = item.Image.Title
		} else if item.ITunesExt != nil && item.ITunesExt.Image != "" {
			image = item.ITunesExt.Image
		}
		if len(item.Categories) > 0 {
			categories = strings.Join(item.Categories, ", ")
		}
		var feedItem models.FeedItem
		err := db.Get(&feedItem, "SELECT * FROM feed_items WHERE guid = ?;", item.GUID)
		if image == "" {
			image, _ = utils.GetOGImage(item.Link)
			if image == "" && feedImage != "" {
				image = feedImage
			}
		}
		if categories == "" && feedCategories != "" {
			categories = feedCategories
		}
		if media == "" && mediaType != "" {
			media = mediaType
		}
		if err != nil {
			pubDate, err := utils.ParseTimeStr(item.Published)
			if err != nil {
				return err
			}
			if _, err = db.Exec(
				"INSERT OR IGNORE INTO feed_items (feed_id, title, link, description, image, alt_text, categories, guid, pub_date, media_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
				feed.FeedID, item.Title, item.Link, item.Description, image, alt, categories, item.GUID, utils.Timestamp(pubDate), media,
			); err != nil {
				return err
			}
		} else {
			if feedItem.Image == "" && image != "" {
				if _, err := db.Exec("UPDATE feed_items SET image = ?, alt_text WHERE id = ?", image, feedItem.ID); err != nil {
					return err
				}
			}
			if feedItem.Categories == "" && categories != "" {
				if _, err := db.Exec("UPDATE feed_items SET categories = ? WHERE id = ?", categories, feedItem.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (h *Handler) FetchRSSFeeds() {
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
				if err := h.FetchRSSFeed(f); err != nil {
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
			h.FetchRSSFeeds()
		}
	}
}
