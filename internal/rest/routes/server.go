package routes

import "github.com/gofiber/fiber/v2"

func (r *Router) CreateServer(c *fiber.Ctx) error {
	return c.SendString("Create Server")
}
