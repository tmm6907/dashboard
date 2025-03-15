package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/tmm6907/dashboard/auth"
	"github.com/tmm6907/dashboard/models"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) GetAllFeeds(c *fiber.Ctx) error {
	var feeds []models.Feed
	db := h.GetDB()
	if err := db.Select(&feeds, "SELECT * FROM feeds;"); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(map[string]any{
		"feeds": feeds,
	})

}

func (h *Handler) GetFeedItems(c *fiber.Ctx) error {
	var feedItems []models.FeedItem
	db := h.GetDB()
	UUIDString := c.Query("feedID")
	if UUIDString != "" {
		feedID, err := h.ParseUUIDString(UUIDString)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		db := h.GetDB()
		if err := db.Select(&feedItems, "SELECT * FROM feed_items WHERE feed_id = ? order by id desc;", feedID); err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}
	if err := db.Select(&feedItems, "SELECT * FROM feed_items;"); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	sort.Slice(feedItems, func(i, j int) bool {
		return time.Time(feedItems[i].PubDate).After(time.Time(feedItems[j].PubDate))
	})
	log.Debug(len(feedItems))

	return c.JSON(map[string]any{
		"feeds": feedItems,
	})
}

func (h *Handler) CreateFeed(c *fiber.Ctx) error {
	db := h.GetDB()
	request := make(map[string]any)
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	log.Debug(request)

	feedUUID := uuid.New()
	feedID := feedUUID[:]

	feedLink := request["link"].(string)

	if utils.IsYoutubeChannelURL(feedLink) {
		link, err := utils.GetYouTubeRSS(feedLink)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		feedLink = link
	}

	if !h.ValidateURL(feedLink) {
		return c.Status(http.StatusBadRequest).SendString("invalid RSS feed link")
	}

	if _, err := db.Exec("INSERT INTO feeds(feed_id, title, link, description, language) VALUES (?, ?, ?, ?, ?);",
		feedID, request["title"], feedLink, request["description"], request["language"]); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return nil
}

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	url := auth.GetLoginURL(h.GetOauthConfig(), "random-state")
	return c.Redirect(url, http.StatusFound)
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
			return c.Status(http.StatusInternalServerError).SendString("Invalid user ID")
		}
		firstName, _ := userInfo["given_name"].(string) // First name
		lastName, _ := userInfo["family_name"].(string) // Last name

		// Check if user exists, if not, create them
		var userID int
		db := h.GetDB()
		err = db.QueryRow("SELECT id FROM users WHERE oauth_id = ? AND oauth_provider = 'google'", oauthID).
			Scan(&userID)

		newUUID := uuid.New()
		mashboardEmail := newUUID.String() + "@mash.board"

		if err == sql.ErrNoRows {
			// Insert new user
			res, err := db.Exec("INSERT INTO users (oauth_provider, oauth_id, first_name, last_name, mashboard_email) VALUES (?, ?, ?, ?, ?)",
				"google", oauthID, firstName, lastName, mashboardEmail)
			if err != nil {
				log.Error("Failed to insert user:", err)
				return c.Status(http.StatusInternalServerError).SendString("User creation failed")
			}
			lastID, err := res.LastInsertId()
			if err == nil {
				userID = int(lastID)
			}
		} else if err != nil {
			log.Error("Database query failed:", err)
			return c.Status(http.StatusInternalServerError).SendString("Database error")
		}

		token, err := h.GenerateJWTToken(oauthID)
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
		return c.SendString("Login successful")
	}
}
