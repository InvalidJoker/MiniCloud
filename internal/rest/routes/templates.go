package routes

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) CreateTemplate(c *fiber.Ctx) error {

	return c.SendString("Create Template")
}
