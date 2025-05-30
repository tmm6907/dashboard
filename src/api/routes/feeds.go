package routes

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/tmm6907/dashboard/models"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) GetFeeds(c *fiber.Ctx) error {
	var feeds []models.Feed
	db := h.GetDB()
	body := struct {
		Query string `json:"query"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(500)
	}
	if body.Query == "" {
		if err := db.Select(&feeds, "SELECT * FROM feeds ORDER BY title;"); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		param := "%" + body.Query + "%"
		if err := db.Select(&feeds, "SELECT * FROM feeds WHERE title LIKE ? OR link LIKE ? OR categories LIKE ? ORDER BY title;", param, param, param); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.JSON(map[string]any{
		"feeds": feeds,
	})
}

// func (h *Handler) CreateFeed(c *fiber.Ctx) error {
// 	db := h.GetDB()
// 	request := make(map[string]any)
// 	if err := json.Unmarshal(c.Body(), &request); err != nil {
// 		return c.Status(http.StatusInternalServerError).SendString(err.Error())
// 	}

// 	link, ok := request["link"]
// 	if !ok {
// 		return c.Status(http.StatusBadRequest).SendString("link is required")
// 	}
// 	title, ok := request["title"]
// 	if !ok {
// 		return c.Status(http.StatusBadRequest).SendString("title is required")
// 	}
// 	feedLink := link.(string)

// 	isYoutube := false
// 	log.Debug("Feed link before youtube check", feedLink)
// 	if utils.IsYoutubeChannelURL(feedLink) {
// 		isYoutube = true
// 		link, err := utils.GetYouTubeRSS(feedLink)
// 		if err != nil {
// 			log.Error(err, link)
// 			return c.Status(http.StatusInternalServerError).SendString(err.Error())
// 		}
// 		feedLink = link
// 	}
// 	log.Debug("Feed link after youtube check", feedLink)

// 	if !h.ValidateURL(feedLink) {
// 		log.Error("Invalid feed url", feedLink)
// 		return c.Status(http.StatusBadRequest).SendString("invalid RSS feed link")
// 	}

// 	description, ok := request["description"]
// 	if !ok {
// 		description = ""
// 	}
// 	language, ok := request["language"]
// 	if !ok {
// 		language = ""
// 	}

// 	feedUUID := uuid.New()
// 	feedID := feedUUID[:]
// 	if isYoutube {
// 		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language, categories, media_type) VALUES (?, ?, ?, ?, ?, ?, ?);",
// 			feedID, title, feedLink, description, language, "youtube", "video"); err != nil {
// 			return c.Status(http.StatusInternalServerError).SendString(err.Error())
// 		}
// 	} else {
// 		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language) VALUES (?, ?, ?, ?, ?);",
// 			feedID, title, feedLink, description, language); err != nil {
// 			return c.Status(http.StatusInternalServerError).SendString(err.Error())
// 		}
// 	}

// 	return nil
// }

func (h *Handler) SearchForNewFeedByURL(c *fiber.Ctx) error {
	body := struct {
		URL string `json:"url"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(500)
	}
	isYoutube := false

	if body.URL != "" {
		if utils.IsYoutubeChannelURL(body.URL) {
			link, err := utils.GetYouTubeRSS(body.URL)
			if err != nil {
				log.Error(err)
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			isYoutube = true
			body.URL = link
		}
		rssParser := gofeed.NewParser()
		feedData, err := rssParser.ParseURL(body.URL)
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		if feedData.Image == nil {
			ogImage, err := utils.GetOGImage(body.URL)
			if err == nil {
				feedData.Image = &gofeed.Image{URL: ogImage}
			} else {
				log.Error(err)
			}
		}
		if isYoutube {
			feedData.Description = fmt.Sprintf("%s YouTube channel", feedData.Title)
			feedData.Title = "YouTube | " + feedData.Title

		}
		data := map[string]any{
			"title":       feedData.Title,
			"description": feedData.Description,
			"image":       feedData.Image,
			"items":       feedData.Items,
		}
		return c.JSON(data)
	}
	return c.Status(http.StatusBadRequest).SendString("url must not be empty")

}

func (h *Handler) FollowFeed(c *fiber.Ctx) error {
	user, err := h.GetUserFromToken(c.Cookies("token"))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	feedData := struct {
		FeedID string `json:"feedID"`
	}{}
	if err = c.BodyParser(&feedData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	var feed *models.Feed
	db := h.GetDB()
	u, _ := uuid.Parse(feedData.FeedID)
	b2 := u[:]
	feedUUID := utils.UUID(b2)
	if err = db.Get(&feed, "SELECT * FROM feeds WHERE feed_id = ?;", feedUUID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString("feed not found")
	}
	if _, err = db.Exec("INSERT INTO feed_follows (user_id, feed_id) VALUES (?, ?);", user.ID, feedUUID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString("failed to follow feed")
	}
	return nil
}

func (h *Handler) CreateAndFollowFeed(c *fiber.Ctx) error {
	user, err := h.GetUserFromToken(c.Cookies("token"))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	feedData := struct {
		Link       string `json:"link"`
		Title      string `json:"title"`
		Desc       string `json:"desc"`
		Collection string `json:"collection"`
	}{}
	if err = c.BodyParser(&feedData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if feedData.Link == "" {
		return c.Status(fiber.StatusBadRequest).SendString("url must not be empty")
	}

	if utils.IsYoutubeChannelURL(feedData.Link) {
		link, err := utils.GetYouTubeRSS(feedData.Link)
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		feedData.Link = link
	}
	db := h.GetDB()
	var feed models.Feed
	if err = db.Get(&feed, "SELECT * FROM feeds WHERE link = ?;", feedData.Link); err != nil {
		feedID := uuid.New()
		if _, err := db.Exec("INSERT OR IGNORE INTO feeds (feed_id, title, link, description) VALUES (?, ?, ?, ?);", utils.UUID(feedID[:]), feedData.Title, feedData.Link, feedData.Desc); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		if err = db.Get(&feed, "SELECT * FROM feeds WHERE feed_id = ?;", utils.UUID(feedID[:])); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	if feedData.Title != "" || feedData.Desc != "" {
		if _, err = db.Exec("INSERT OR IGNORE INTO feed_follows (user_id, feed_id, user_feed_name, user_feed_desc) VALUES (?, ?, ?, ?);", user.ID, feed.FeedID, feedData.Title, feedData.Desc); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if _, err = db.Exec("INSERT OR IGNORE INTO feed_follows (user_id, feed_id) VALUES (?, ?);", user.ID, feedData.Link); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	if feedData.Collection != "" {
		var collection_id int
		if err := db.Get(&collection_id, "SELECT id from collections WHERE name = ?", feedData.Collection); err != nil {
			if _, err = db.Exec("INSERT INTO collections (name) VALUES (?);", feedData.Collection); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			if err := db.Get(&collection_id, "SELECT id from collections WHERE name = ?", feedData.Collection); err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
		}
		if _, err := db.Exec("INSERT OR IGNORE INTO user_collections (user_id, collection_id) VALUES (?, ?);", user.ID, collection_id); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) GetFollowedFeeds(c *fiber.Ctx) error {
	var feeds []models.Feed
	db := h.GetDB()
	body := struct {
		Query string `json:"query"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		log.Error(err)
		return c.SendStatus(500)
	}
	user, err := h.GetUserFromToken(c.Cookies("token"))
	if err != nil {
		log.Error(err)
		return c.SendStatus(http.StatusUnauthorized)
	}
	if body.Query == "" {
		if err := db.Select(&feeds, "SELECT feeds.* FROM feeds JOIN feed_follows ff ON feeds.feed_id = ff.feed_id WHERE ff.user_id = ? ORDER BY title;", user.ID); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		param := "%" + body.Query + "%"
		if err := db.Select(&feeds, "SELECT feeds.* FROM feeds JOIN feed_follows ff ON feeds.feed_id = ff.feed_id WHERE ff.user_id = ? AND (title LIKE ? OR link LIKE ? OR categories LIKE ?) ORDER BY title;", user.ID, param, param, param); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.JSON(map[string]any{
		"feeds": feeds,
	})
}
