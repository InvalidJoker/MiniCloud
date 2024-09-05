package routes

import (
	"minicloud/internal/cloud"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	DockerService *cloud.DockerService
	Fiber         *fiber.App
}

func NewRouter(dockerService *cloud.DockerService, app *fiber.App) *Router {
	return &Router{
		DockerService: dockerService,
		Fiber:         app,
	}
}
