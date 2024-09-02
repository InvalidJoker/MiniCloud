package cloud

import (
	"context"
	"minicloud/internal/config"
	"minicloud/internal/database"

	"github.com/docker/docker/client"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type DockerService struct {
	Client   *client.Client
	Context  context.Context
	Proxy    *proxy.Proxy
	Database *database.Database
	Config   config.Config
}

func NewDockerService(
	db *database.Database,
	config config.Config,
	proxy *proxy.Proxy,
) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerService{
		Client:   cli,
		Context:  context.Background(),
		Database: db,
		Config:   config,
		Proxy:    proxy,
	}, nil
}

func (s *DockerService) DockerTest(ctx *context.Context) error {
	server := &database.Server{
		Name: "test",
		Port: 25566,
	}

	_, err := s.CreateServer(*ctx, server)
	if err != nil {
		return err
	}

	return nil

}
