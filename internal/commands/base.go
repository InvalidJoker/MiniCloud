package commands

import (
	"minicloud/internal/cloud"
	"minicloud/internal/database"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type CommandService struct {
	Database      *database.Database
	DockerService *cloud.DockerService
	Proxy         *proxy.Proxy
}

func NewCommandService(db *database.Database, dockerService *cloud.DockerService, proxy *proxy.Proxy) *CommandService {
	return &CommandService{
		Database:      db,
		DockerService: dockerService,
		Proxy:         proxy,
	}
}

func (s *CommandService) Register() {}
