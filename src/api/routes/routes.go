package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/tmm6907/dashboard/auth"
	"github.com/tmm6907/dashboard/models"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) GetFeeds(c *fiber.Ctx) error {
	var feeds []models.Feed
	db := h.GetDB()
	query := c.Query("query")
	if query == "" {
		if err := db.Select(&feeds, "SELECT * FROM feeds ORDER BY title;"); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if err := db.Select(&feeds, "SELECT * FROM feeds WHERE title LIKE ? OR link LIKE ? ORDER BY title;", "%"+query+"%", "%"+query+"%"); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.JSON(map[string]any{
		"feeds": feeds,
	})
}

func (h *Handler) GetFeedItem(c *fiber.Ctx) error {
	feedItemID := c.Params("id")
	feedItem := make(map[string]interface{})

	db := h.GetDB()
	row := db.QueryRowx(`
		SELECT fi.*, f.title as feedName 
		FROM feed_items fi
		JOIN feeds f ON fi.feed_id = f.feed_id
		WHERE fi.id = ?;`, feedItemID)

	if err := row.MapScan(feedItem); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	log.Debug(feedItem["pub_date"])
	return c.JSON(feedItem)
}

func (h *Handler) GetFeedItems(c *fiber.Ctx) error {
	var feedItems []models.FeedItem
	db := h.GetDB()
	token := c.Cookies("token")
	if token == "" {
		log.Error("token should not be empty")
		return c.SendStatus(http.StatusInternalServerError)
	}
	user, err := h.GetUserFromToken(token)
	if err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	userID := user.ID

	category := strings.ToLower(c.Query("category"))
	if category != "" && category != "all" && category != "all categories" {
		if category == "technology" {
			category = "tech"
		}
		log.Debug(category)
		if err := db.Select(&feedItems, "SELECT * FROM feed_items WHERE categories LIKE ? OR media_type LIKE ?;", "%"+category+"%", "%"+category+"%"); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if err := db.Select(&feedItems, "SELECT * FROM feed_items;"); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	sort.Slice(feedItems, func(i, j int) bool {
		return time.Time(feedItems[i].PubDate).After(time.Time(feedItems[j].PubDate))
	})

	latest := []map[string]any{}
	var saved []models.FeedItem
	var collections []string

	if err = db.Select(&saved, `
			SELECT fi.* 
			FROM feed_items fi
			JOIN saved_feeds sf ON fi.id = sf.feed_item_id
			WHERE sf.user_id = ?;`, userID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if err = db.Select(&collections, `
			SELECT c.name 
			FROM collections c
			JOIN user_collections uc ON c.id = uc.collection_id
			WHERE uc.user_id = ?;`, userID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	items := []map[string]any{}
	for _, item := range feedItems {
		var name string
		err = db.Get(&name, "SELECT title from feeds where feed_id = ?", item.FeedID)
		if err != nil {
			log.Error(err)
		}
		resultItem := make(map[string]any)
		resultItem["name"] = name
		resultItem["item"] = item
		pubDate := time.Time(item.PubDate)
		if time.Now().Sub(pubDate) <= 3*24*time.Hour {
			latest = append(latest, resultItem)
		}
		items = append(items, resultItem)
	}

	log.Debug(len(latest), len(saved), len(collections), len(feedItems))
	if len(latest) > 0 {
		log.Debug(latest[0]["name"])
	}
	return c.JSON(map[string]any{
		"latest":      latest,
		"items":       items,
		"saved":       saved,
		"collections": collections,
	})
}

func (h *Handler) CreateFeed(c *fiber.Ctx) error {
	db := h.GetDB()
	request := make(map[string]any)
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	link, ok := request["link"]
	if !ok {
		return c.Status(http.StatusBadRequest).SendString("link is required")
	}
	title, ok := request["title"]
	if !ok {
		return c.Status(http.StatusBadRequest).SendString("title is required")
	}
	feedLink := link.(string)

	isYoutube := false
	if utils.IsYoutubeChannelURL(feedLink) {
		isYoutube = true
		link, err := utils.GetYouTubeRSS(feedLink)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		feedLink = link
	}

	if !h.ValidateURL(feedLink) {
		return c.Status(http.StatusBadRequest).SendString("invalid RSS feed link")
	}

	description, ok := request["description"]
	if !ok {
		description = ""
	}
	language, ok := request["language"]
	if !ok {
		language = ""
	}

	feedUUID := uuid.New()
	feedID := feedUUID[:]
	if isYoutube {
		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language, categories, media_type) VALUES (?, ?, ?, ?, ?, ?, ?);",
			feedID, title, feedLink, description, language, "youtube", "video"); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language) VALUES (?, ?, ?, ?, ?);",
			feedID, title, feedLink, description, language); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	return nil
}

func (h *Handler) SearchForFeed(c *fiber.Ctx) error {
	feedURL := c.Query("url")
	rssParser := gofeed.NewParser()
	feedData, err := rssParser.ParseURL(feedURL)
	if err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if feedData.Image == nil {
		ogImage, err := utils.GetOGImage(feedURL)
		if err == nil {
			feedData.Image = &gofeed.Image{URL: ogImage}
		} else {
			log.Error(err)
		}
	}

	data := map[string]any{
		"title":       feedData.Title,
		"description": feedData.Description,
		"image":       feedData.Image,
		"items":       feedData.Items,
	}
	return c.JSON(data)
}

func (h *Handler) FollowFeed(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("auth token empty")
	}
	user, err := h.GetUserFromToken(token)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	feedData := struct {
		Link       string `json:"link"`
		Title      string `json:"title"`
		Desc       string `json:"desc"`
		Collection string `json:"collection"`
	}{}
	if err = c.BodyParser(&feedData); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if feedData.Link == "" {
		log.Debug("empty link")
		return c.Status(fiber.StatusBadRequest).SendString("url must not be empty")
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
		log.Debug(feed)
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

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	url := auth.GetLoginURL(h.GetOauthConfig(), "random-state")
	log.Debug("Redirecting to ", url)
	return c.Status(fiber.StatusFound).SendString(url)
}

func (h *Handler) CallbackHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Debug("callback reached")
		code := c.Query("code")
		if code == "" {
			log.Error("Google did not return code to callback func")
			return c.Status(http.StatusInternalServerError).SendString("Google auth error")
		}
		userInfo, err := h.FetchGoogleInfo(code)
		if err != nil {
			log.Error("Failed to log into Google OAuth: ", err)
			return c.Status(http.StatusInternalServerError).SendString("Invalid user ID")
		}

		// Check if user exists, if not, create them
		var userID utils.UUID
		db := h.GetDB()
		newUUID := uuid.New()
		mashboardEmail := newUUID.String() + "@mash.board"
		err = db.QueryRow("SELECT id FROM users WHERE oauth_id = ? AND oauth_provider = 'google'", userInfo["oauthID"]).
			Scan(&userID)

		if err == sql.ErrNoRows {
			// Insert new user
			_, err := db.Exec("INSERT INTO users (id, oauth_provider, oauth_id, first_name, last_name, mashboard_email) VALUES (?, ?, ?, ?, ?, ?)",
				newUUID[:], "google", userInfo["oauthID"], userInfo["firstName"], userInfo["lastName"], mashboardEmail)
			if err != nil {
				log.Error("Failed to insert user:", err)
				return c.Status(http.StatusInternalServerError).SendString("User creation failed")
			}
			log.Debug("User created")
		} else if err != nil {
			log.Error("Database query failed:", err)
			return c.Status(http.StatusInternalServerError).SendString("Database error")
		}

		if err = h.SaveNewToken(c, userInfo["oauthID"], userInfo["firstName"], userInfo["lastName"], code); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("Unable to generate auth token")
		}
		return c.Redirect("http://localhost:3030")
	}
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration in the past
		HTTPOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: "Lax",
	})
	return c.Status(fiber.StatusOK).SendString("Logged out")
}
