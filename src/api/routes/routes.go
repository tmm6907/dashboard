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

func (h *Handler) GetFeedItems(c *fiber.Ctx) error {
	var feedItems []models.FeedItem
	db := h.GetDB()
	token := c.Cookies("token")
	if token == "" {
		log.Error("not logged in")
		return c.Redirect("http://localhost:8080/auth/login")
	}
	userID, err := h.GetUserIDFromToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return c.Redirect("http://localhost:8080/auth/login")
		}
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	UUIDString := c.Query("feedID")
	category := strings.ToLower(c.Query("category"))
	if UUIDString != "" {
		feedID, err := h.ParseUUIDString(UUIDString)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		db := h.GetDB()
		if err := db.Select(&feedItems, "SELECT * FROM feed_items WHERE feed_id = ? ;", feedID); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	if category != "" && category != "all" && category != "all categories" {
		log.Debug(category)
		if err := db.Select(&feedItems, "SELECT * FROM feed_items WHERE categories LIKE ?;", "%"+category+"%"); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if err := db.Select(&feedItems, "SELECT * FROM feed_items;"); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	log.Debug(len(feedItems))

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
	token := c.Cookies("token")
	if token == "" {
		log.Debug("not logged in")
	}
	db := h.GetDB()
	request := make(map[string]any)
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	log.Debug(request)

	feedUUID := uuid.New()
	feedID := feedUUID[:]

	feedLink := request["link"].(string)
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
	if isYoutube {
		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language, categories, media_type) VALUES (?, ?, ?, ?, ?, ?, ?);",
			feedID, request["title"], feedLink, request["description"], request["language"], "youtube", "video"); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language) VALUES (?, ?, ?, ?, ?);",
			feedID, request["title"], feedLink, request["description"], request["language"]); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	return nil
}

func (h *Handler) CheckAuth(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		log.Debug("not logged in")
		return c.Redirect("http://localhost:8080/auth/login")
	}
	if _, err := h.GetUserIDFromToken(token); err == sql.ErrNoRows {
		db := h.GetDB()
		newUUID := uuid.New()
		oauthID, firstName, lastName, err := h.GetOauthInfoFromToken(token)
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		mashboardEmail := newUUID.String() + "@mash.board"
		db.Exec("INSERT INTO users (id, oauth_provider, oauth_id, first_name, last_name, mashboard_email) VALUES (?, ?, ?, ?, ?, ?)",
			newUUID[:], "google", oauthID, firstName, lastName, mashboardEmail)
	}
	log.Debug("logged in")
	return c.Status(http.StatusOK).SendString("logged in!")
}

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	url := auth.GetLoginURL(h.GetOauthConfig(), "random-state")
	log.Debug("Redirecting to ", url)
	return c.Status(fiber.StatusFound).SendString(url)
}

func (h *Handler) CallbackHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code", "")
		if code == "" {
			return c.Status(http.StatusBadRequest).SendString("Missing code")
		}

		// Fetch user info from Google
		userInfo, err := auth.GetUserInfo(h.GetOauthConfig(), code)
		if err != nil {
			log.Error("Failed to get user info:", err)
			return c.Status(http.StatusInternalServerError).SendString("OAuth failed")
		}

		// Extract relevant details
		oauthID, ok := userInfo["id"].(string)
		if !ok {
			log.Error("Invalid user id")
			return c.Status(http.StatusInternalServerError).SendString("Invalid user ID")
		}
		firstName, _ := userInfo["given_name"].(string) // First name
		lastName, _ := userInfo["family_name"].(string) // Last name

		// Check if user exists, if not, create them
		var userID utils.UUID
		db := h.GetDB()
		newUUID := uuid.New()
		mashboardEmail := newUUID.String() + "@mash.board"
		err = db.QueryRow("SELECT id FROM users WHERE oauth_id = ? AND oauth_provider = 'google'", oauthID).
			Scan(&userID)

		log.Debug("checking user")

		if err == sql.ErrNoRows {
			// Insert new user
			_, err := db.Exec("INSERT INTO users (id, oauth_provider, oauth_id, first_name, last_name, mashboard_email) VALUES (?, ?, ?, ?, ?, ?)",
				newUUID[:], "google", oauthID, firstName, lastName, mashboardEmail)
			if err != nil {
				log.Error("Failed to insert user:", err)
				return c.Status(http.StatusInternalServerError).SendString("User creation failed")
			}
			log.Debug("User created")
		} else if err != nil {
			log.Error("Database query failed:", err)
			return c.Status(http.StatusInternalServerError).SendString("Database error")
		}

		token, err := h.GenerateJWTToken(oauthID, firstName, lastName)
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("Unable to generate auth token")
		}

		cookie := &fiber.Cookie{
			Name:     "token",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		}

		if c.Protocol() == "https" {
			cookie.Secure = true
		}
		c.Cookie(cookie)
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
