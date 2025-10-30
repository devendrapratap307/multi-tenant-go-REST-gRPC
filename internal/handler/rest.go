package handler

import (
	"go-multitenant/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RestHandler struct {
	UserService *service.UserService
}

func NewRestHandler(us *service.UserService) *RestHandler { return &RestHandler{UserService: us} }
func (h *RestHandler) Register(app *fiber.App) {
	app.Get("/users/:id", h.GetUser)
}
func (h *RestHandler) GetUser(c *fiber.Ctx) error {
	clientID := c.Get("X-Client-ID")
	if clientID == "" {
		return c.Status(400).SendString("missing X-Client-ID header")
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("invalid id")
	}
	u, err := h.UserService.GetUser(c.Context(), clientID, uint(id))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(u)
}
