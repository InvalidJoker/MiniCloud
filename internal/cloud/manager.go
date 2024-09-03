package cloud

import (
	"context"
	"fmt"
	"minicloud/internal/database"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (s *DockerService) CreateServer(ctx context.Context, server *database.Server) (DockerServer, error) {

	// check if container already exists
	_, err := s.Client.ContainerInspect(ctx, server.Name)
	if err == nil {
		return DockerServer{}, fmt.Errorf("container already exists")
	}

	env := []string{
		"EULA=TRUE",
		"ENABLE_RCON=false",
		"ONLINE_MODE=false",
		"USE_NATIVE_TRANSPORT=false",
	}

	if server.Template.Software != "" {
		env = append(env, "TYPE="+server.Template.Software)
	} else {
		env = append(env, "TYPE=PAPER")
	}

	if server.Version != "" {
		env = append(env, "VERSION="+server.Version)
	}

	CreateServer(server.Name)

	// move template to server
	server.Template.MoveToServer(server.Name)
	// save server data in /data/servers/servername

	fmt.Printf("Source Path: %s\n", filepath.Join("data", "servers", server.Name))
	fmt.Printf("Server Name: %s\n", server.Name)

	sourcePath, err := filepath.Abs(filepath.Join("data", "servers", server.Name))
	if err != nil {
		return DockerServer{}, err
	}

	resp, err := s.Client.ContainerCreate(ctx, &container.Config{
		Image: "itzg/minecraft-server",
		ExposedPorts: nat.PortSet{
			"25565/tcp": struct{}{},
		},
		Env:         env,
		Tty:         true,
		AttachStdin: true,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(server.Port),
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourcePath,
				Target: "/data",
			},
		},
	}, nil, nil, server.Name)
	if err != nil {
		return DockerServer{}, err
	}

	if err := s.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return DockerServer{}, err
	}

	if err := s.RegisterServer(ctx, server); err != nil {
		return DockerServer{}, err
	}

	server.ID = resp.ID

	s.Database.Save(server)

	st, err := s.Client.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return DockerServer{}, err
	}

	return DockerServer{
		Server: server,
		Stream: &st,
	}, nil

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
				if strings.Contains(err.Error(), "container already started") {
					continue
				}
				if strings.Contains(err.Error(), "No such container") {
					if _, err := s.CreateServer(ctx, &server); err != nil {
						return err
					}
				}
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
