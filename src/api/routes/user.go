package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetUser(c *fiber.Ctx) error {
	token := c.Cookies("token")
	user, err := h.GetUserFromToken(token)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}
	return c.JSON(user)
}
