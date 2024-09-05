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

func (s *DockerService) CreateServer(ctx context.Context, req *database.CreateServerRequest) (DockerServer, error) {

	var template database.Template
	s.Database.Where("name = ?", req.Template).First(&template)

	if template.Name == "" {
		// create default template
		template = database.Template{
			Name:     req.Template,
			Software: "paper",
			Version:  "latest",
		}

		s.Database.Create(&template)
		CreateTemplate(template.Name)

	}

	server := &database.Server{
		Name:     req.Name,
		Port:     req.Port,
		Lobby:    req.Lobby,
		Template: template,
	}

	if req.CustomData != nil {
		server.CustomData = req.CustomData
	}

	// check if container already exists
	_, err := s.Client.ContainerInspect(ctx, server.Name)
	if err == nil {
		return DockerServer{}, nil
	}

	var env []string
	var image string

	if template.CustomImage == "" {
		image = "itzg/minecraft-server"
		env = append(env, []string{
			"EULA=TRUE",
			"ENABLE_RCON=false",
			"ONLINE_MODE=false",
			"USE_NATIVE_TRANSPORT=false",
		}...)

		if server.Template.Software != "" {
			env = append(env, "TYPE="+server.Template.Software)
		} else {
			env = append(env, "TYPE=PAPER")
		}

		if server.Template.Version != "" {
			env = append(env, "VERSION="+server.Template.Version)
		}
	} else {
		// use custom image
		image = server.Template.CustomImage

		if server.Template.CustomImageData != nil {
			for key, value := range server.Template.CustomImageData {
				env = append(env, fmt.Sprintf("%s=%s", strconv.Itoa(key), string(value)))
			}

			for key, value := range server.CustomData {
				env = append(env, fmt.Sprintf("%s=%s", strconv.Itoa(key), string(value)))
			}
		}
	}

	_, err = CreateServer(server.Name)

	// move template to server
	err = server.Template.MoveToServer(server.Name)

	if err != nil {
		return DockerServer{}, err
	}

	fmt.Printf("Source Path: %s\n", filepath.Join("data", "servers", server.Name))
	fmt.Printf("Server Name: %s\n", server.Name)

	sourcePath, err := filepath.Abs(filepath.Join("data", "servers", server.Name))
	if err != nil {
		return DockerServer{}, err
	}

	resp, err := s.Client.ContainerCreate(ctx, &container.Config{
		Image: image,
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

func (s *DockerService) GetServerStatus(ctx context.Context, server *database.Server) (int, error) {
	c, err := s.Client.ContainerInspect(ctx, server.ID)

	if err != nil {
		return -1, err
	}

	if c.State.Running {
		return 1, nil
	}

	if c.State.Restarting {
		return 2, nil
	}

	return -1, nil
}

func (s *DockerService) LoadServers(ctx context.Context) error {

	var servers []database.Server
	s.Database.Find(&servers)

	fmt.Printf("Servers: %v\n", servers)

	for _, server := range servers {
		fmt.Printf("Server: %s\n", server.Name)
		if server.ID != "" {
			if err := s.StartServer(ctx, &server); err != nil {
				if strings.Contains(err.Error(), "container already started") {
					continue
				}
				if strings.Contains(err.Error(), "No such container") {

					err = server.Template.MoveToServer(server.Name)

					if err != nil {
						return err
					}
					if _, err := s.CreateServer(ctx, server.ToRequest()); err != nil {
						return err
					}
				}
			}
		} else {
			if _, err := s.CreateServer(ctx, server.ToRequest()); err != nil {
				return err
			}
		}

		_, err := s.Proxy.Register(server.GetServerInfo())

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DockerService) DeleteServer(ctx context.Context, server *database.Server) error {
	if err := s.Client.ContainerStop(ctx, server.ID, container.StopOptions{}); err != nil {
		return err
	}

	if err := s.Client.ContainerRemove(ctx, server.ID, container.RemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func (s *DockerService) StopServer(ctx context.Context, server *database.Server) error {
	return s.Client.ContainerStop(ctx, server.ID, container.StopOptions{})
}

func (s *DockerService) ToDockerServer(server *database.Server) DockerServer {

	st, err := s.Client.ContainerAttach(s.Context, server.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})

	if err != nil {
		return DockerServer{
			Client: s.Client,
			Server: server,
		}
	}

	return DockerServer{
		Client: s.Client,
		Server: server,
		Stream: &st,
	}
}
