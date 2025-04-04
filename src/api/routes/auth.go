package routes

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/tmm6907/dashboard/auth"
	"github.com/tmm6907/dashboard/utils"
)

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	url := auth.GetLoginURL(h.GetOauthConfig(), "random-state")
	return c.Status(fiber.StatusFound).SendString(url)
}

func (h *Handler) CallbackHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
		} else if err != nil {
			log.Error("Database query failed:", err)
			return c.Status(http.StatusInternalServerError).SendString("Database error")
		}

		if err = h.SaveNewToken(c, userInfo["oauthID"], userInfo["firstName"], userInfo["lastName"], code); err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("Unable to generate auth token")
		}
		return c.Redirect("http://localhost:4173")
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
