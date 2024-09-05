package rest

import (
	"minicloud/internal/cloud"
	"minicloud/internal/config"
	"minicloud/internal/rest/routes"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type BackendService struct {
	DockerService *cloud.DockerService
	Config        config.Config
}

func NewBackendService(dockerService *cloud.DockerService, config config.Config) *BackendService {
	return &BackendService{
		DockerService: dockerService,
		Config:        config,
	}
}

func (b *BackendService) Start() {
	app := fiber.New()

	app.Post("/start", b.start)

	router := routes.NewRouter(b.DockerService, app)

	if b.Config.AuthToken != "" {
		app.Use(func(c *fiber.Ctx) error {
			// check if the request is authorized
			authHeader := strings.Split(c.Get("Authorization"), "Bearer ")

			if len(authHeader) != 2 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			// check if the token is valid
			if authHeader[1] != b.Config.AuthToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			return c.Next()
		})
	}

	router.Fiber = app

	err := app.Listen(strconv.Itoa(b.Config.Port))

	if err != nil {
		panic(err)
	}

}

func (b *BackendService) start(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
