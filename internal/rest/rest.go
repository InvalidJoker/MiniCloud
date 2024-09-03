package rest

import (
	"minicloud/internal/cloud"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BackendService struct {
	DockerService *cloud.DockerService
	Fiber         *fiber.App
}

func NewBackendService(dockerService *cloud.DockerService) *BackendService {
	return &BackendService{
		DockerService: dockerService,
	}
}

func (b *BackendService) Start() {
	app := fiber.New()

	app.Post("/start", b.start)

	b.Fiber = app

	http.ListenAndServe(":8080", nil)

}

func (b *BackendService) start(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
