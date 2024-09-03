package rest

import (
	"minicloud/internal/cloud"
	"minicloud/internal/rest/routes"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BackendService struct {
	DockerService *cloud.DockerService
	Router        *routes.Router
}

func NewBackendService(dockerService *cloud.DockerService) *BackendService {
	return &BackendService{
		DockerService: dockerService,
	}
}

func (b *BackendService) Start() {
	app := fiber.New()

	app.Post("/start", b.start)

	b.Router.Fiber = app

	http.ListenAndServe(":8080", nil)

}

func (b *BackendService) start(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
