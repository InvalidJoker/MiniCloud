package routes

import (
	"minicloud/internal/database"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) CreateServer(c *fiber.Ctx) error {

	// Initialize the req variable properly
	req := new(database.CreateServerRequest)

	// Parse body
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create server
	server, err := r.DockerService.CreateServer(r.DockerService.Context, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(server)

}
