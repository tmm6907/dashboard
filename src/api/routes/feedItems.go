package routes

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tmm6907/dashboard/models"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) GetFeedItems(c *fiber.Ctx) error {
	feedItems := []map[string]any{}
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
		err = h.QueryRows(&feedItems, "SELECT fi.*, f.title as feed_name FROM feed_items fi JOIN feeds f ON fi.feed_id = f.feed_id WHERE categories LIKE ? OR media_type LIKE ?;", "%"+category+"%", "%"+category+"%")
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		err = h.QueryRows(&feedItems, "SELECT fi.*, f.title as feed_name FROM feed_items fi JOIN feeds f ON fi.feed_id = f.feed_id;")
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	sort.Slice(feedItems, func(i, j int) bool {
		iDateStr := feedItems[i]["pub_date"].(string)
		iPubdate, _ := utils.ParseTimeStr(iDateStr)
		jDateStr := feedItems[i]["pub_date"].(string)
		jPubDate, _ := utils.ParseTimeStr(jDateStr)
		return iPubdate.After(jPubDate)
	})

	latest := []map[string]any{}
	// var saved []models.FeedItem
	var collections []models.Collection

	// if err = db.Select(&saved, `
	// 		SELECT fi.*
	// 		FROM feed_items fi
	// 		JOIN saved_feeds sf ON fi.id = sf.feed_item_id
	// 		WHERE sf.user_id = ?;`, userID); err != nil {
	// 	log.Error(err)
	// 	return c.Status(http.StatusInternalServerError).SendString(err.Error())
	// }

	if err = db.Select(&collections, `
			SELECT c.*
			FROM collections c
			JOIN user_collections uc ON c.id = uc.collection_id
			WHERE uc.user_id = ?;`, userID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	for _, item := range feedItems {
		dateStr := item["pub_date"].(string)
		pubDate, err := utils.ParseTimeStr(dateStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		if time.Now().Sub(pubDate) <= 3*24*time.Hour {
			latest = append(latest, item)
		}
	}

	log.Debug(len(latest), len(collections), len(feedItems))
	return c.JSON(map[string]any{
		"latest":      latest,
		"items":       feedItems,
		"collections": collections,
	})
}

func (h *Handler) GetFeedItem(c *fiber.Ctx) error {
	token := c.Cookies("token")
	feedItemID := c.Params("id")
	feedItem := make(map[string]interface{})
	log.Debug(c.OriginalURL())
	user, err := h.GetUserFromToken(token)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	db := h.GetDB()
	row := db.QueryRowx(`
		SELECT fi.*, 
           f.title AS feed_name, 
           CASE 
               WHEN sf.feed_item_id IS NOT NULL THEN true 
               ELSE false 
           END AS bookmarked
    FROM feed_items fi
    JOIN feeds f ON fi.feed_id = f.feed_id
    LEFT JOIN saved_feeds sf ON fi.id = sf.feed_item_id AND sf.user_id = ?
    WHERE fi.id = ?;`, user.ID, feedItemID)

	if err := row.MapScan(feedItem); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(feedItem)
}

func (h *Handler) SaveFeedItem(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Status(http.StatusInternalServerError).SendString("expected auth token")
	}
	user, err := h.GetUserFromToken(token)
	if err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString("unable to determine user id")
	}
	feedItemID := c.Params("id")

	db := h.GetDB()
	var savedFeed models.SavedFeed
	if err = db.Get(&savedFeed, "SELECT * FROM saved_feeds WHERE user_id = ? AND feed_item_id = ? ;", user.ID, feedItemID); err == nil {
		log.Debug("removing")
		if _, err := db.Exec("DELETE FROM saved_feeds where user_id = ? and feed_item_id = ? ;", user.ID, feedItemID); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("unexpected error")
		}
		return c.SendStatus(200)
	}
	if _, err = db.Exec("INSERT INTO saved_feeds (user_id, feed_item_id) VALUES (?, ?);", user.ID, feedItemID); err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString("unable to save feed item")
	}
	return c.SendStatus(200)
}

func (h *Handler) GetSavedFeedItems(c *fiber.Ctx) error {
	log.Debug("called")
	token := c.Cookies("token")
	if token == "" {
		return c.Status(http.StatusInternalServerError).SendString("expected auth token")
	}
	user, err := h.GetUserFromToken(token)
	if err != nil {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).SendString("unable to determine user id")
	}
	items := []map[string]any{}
	err = h.QueryRows(&items, "SELECT fi.*, f.title as feed_name FROM feed_items fi JOIN saved_feeds sf ON fi.id = sf.feed_item_id LEFT JOIN feeds f ON fi.feed_id = f.feed_id WHERE sf.user_id = ?;", user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("unable to determine user id")
	}

	return c.JSON(map[string]any{
		"items": items,
	})
}
