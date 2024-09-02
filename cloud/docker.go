package cloud

import (
	"context"
	"minicloud/config"
	"minicloud/database"

	"github.com/docker/docker/client"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type DockerService struct {
	Client  *client.Client
	Context context.Context
	Proxy   *proxy.Proxy
}

func NewDockerService(
	db *database.Database,
	config config.Config,
) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerService{
		Client:  cli,
		Context: context.Background(),
	}, nil
}
