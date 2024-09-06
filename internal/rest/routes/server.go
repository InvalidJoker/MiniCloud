package routes

import (
	"minicloud/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func (r *Router) CreateServer(c *fiber.Ctx) error {

	//database.CreateServerRequest{}

	req := &database.CreateServerRequest{
		Name:     c.Query("name"),
		Port:     c.QueryInt("port"),
		Lobby:    c.QueryBool("lobby"),
		Template: c.Query("template"),

		CustomData: datatypes.JSON([]byte(c.Query("custom_data"))),
	}

	validate := validator.New()
	err := validate.Struct(req)
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
