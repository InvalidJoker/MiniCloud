package cloud

import (
	"context"
	"minicloud/internal/database"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func (s *DockerService) CreateServer(ctx context.Context, server *database.Server) (string, error) {

	// check if container already exists
	_, err := s.Client.ContainerInspect(ctx, server.Name)
	if err == nil {
		return "", nil
	}

	env := []string{
		"EULA=TRUE",
		"ENABLE_RCON=false",
	}

	if server.Software != "" {
		env = append(env, "TYPE="+server.Software)
	}

	if server.Version != "" {
		env = append(env, "VERSION="+server.Version)
	}

	resp, err := s.Client.ContainerCreate(ctx, &container.Config{
		Image: "itzg/minecraft-server",
		ExposedPorts: nat.PortSet{
			"25565/tcp": struct{}{},
		},
		Env: env,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(server.Port),
				},
			},
		},
	}, nil, nil, server.Name)
	if err != nil {
		return "", err
	}

	if err := s.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	if err := s.RegisterServer(ctx, server); err != nil {
		return "", err
	}

	server.ID = resp.ID

	s.Database.Create(server)

	return resp.ID, nil

}

func (s *DockerService) StartServer(ctx context.Context, server *database.Server) error {
	return s.Client.ContainerStart(ctx, server.ID, container.StartOptions{})
}

func (s *DockerService) LoadServers(ctx context.Context) error {

	var servers []database.Server
	s.Database.Find(&servers)

	for _, server := range servers {
		if server.ID != "" {
			if err := s.StartServer(ctx, &server); err != nil {
				return err
			}
		} else {
			if _, err := s.CreateServer(ctx, &server); err != nil {
				return err
			}
		}

		s.Proxy.Register(server.GetServerInfo())
	}

	return nil
}
