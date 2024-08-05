package cloud

import (
	"context"
	"github.com/docker/docker/client"
)


type DockerService struct {
	Client *client.Client
	Context context.Context
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerService{
		Client: cli,
		Context: context.Background(),
	}, nil
}
