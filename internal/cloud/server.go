package cloud

import (
	"context"
	"minicloud/internal/database"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerServer struct {
	Server *database.Server
	Client *client.Client
	Stream *types.HijackedResponse
}

func (s *DockerServer) Online() bool {
	c, err := s.Client.ContainerInspect(context.Background(), s.Server.ID)

	if err != nil {
		return false
	}

	if c.State.Running {
		return true
	}

	return false

}

func (s *DockerServer) Start() error {
	return s.Client.ContainerStart(context.Background(), s.Server.ID, container.StartOptions{})
}

func (s *DockerServer) Stop() error {
	return s.Client.ContainerStop(context.Background(), s.Server.ID, container.StopOptions{})
}

func (s *DockerServer) Restart() error {
	return s.Client.ContainerRestart(context.Background(), s.Server.ID, container.StopOptions{})
}
