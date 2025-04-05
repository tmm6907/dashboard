package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetUser(c *fiber.Ctx) error {
	token := c.Cookies("token")
	user, err := h.GetUserFromToken(token)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)
}
