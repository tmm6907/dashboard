package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetUser(c *fiber.Ctx) error {
	user, err := h.GetUserFromToken(c.Cookies("token"))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	return c.JSON(user)
}
