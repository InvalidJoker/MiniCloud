package cloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	Client  *client.Client
	Context context.Context
}

type ServerSettings struct {
	Ram   int
	Ports []int
}

type DockerContainer struct {
	ID    string
	Image string
	Name  string
	Type  string

	Settings ServerSettings
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerService{
		Client:  cli,
		Context: context.Background(),
	}, nil
}

// Create a new Docker Container based on this configuration
func (ds *DockerService) Create(c *DockerContainer) error {
	ctx := context.Background()

	containerConfig := parseContainerConfig(c)
	containerHostConfig := parseHostConfig(c)

	_, err := ds.Client.ContainerCreate(
		ctx,
		containerConfig,
		containerHostConfig,
		nil,
		nil,
		containerConfig.Hostname,
	)
	if err != nil {
		return err
	}

	//c.server.Save()
	return nil
}

func parseContainerConfig(c *DockerContainer) *container.Config {
	portSet := nat.PortSet{}
	for _, p := range c.Settings.Ports {
		port := nat.Port(fmt.Sprintf("%d/%s", p, "tcp"))
		portSet[port] = struct{}{}
	}

	containerConfig := &container.Config{
		Image:        c.Image,
		AttachStdin:  true,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Hostname:     "daemon-" + c.ID,
		ExposedPorts: portSet,
		Volumes: map[string]struct{}{
			"/data": {},
		},
		// TODO set ram on env variables
		Env: []string{
			"EULA=TRUE",
			"PAPER_DOWNLOAD_URL=https://heroslender.com/assets/PaperSpigot-1.8.8.jar",
			"TYPE=PAPER",
			"VERSION=1.8.8",
			"ENABLE_RCON=false",
		},
	}

	return containerConfig
}

func parseHostConfig(c *DockerContainer) *container.HostConfig {
	portMap := nat.PortMap{}
	for _, p := range c.Settings.Ports {
		port := nat.Port(fmt.Sprintf("%d/%s", p, "tcp"))
		portMap[port] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", p)}}
	}

	// fix windows path
	path := strings.Replace(c.DataPath(), "C:\\", "/c/", 1)
	path = strings.Replace(path, "\\", "/", -1)
	// point to `/data` volume
	path += ":/data"

	memory := c.Settings.Ram * 1024 * 1024 // convert to bytes
	containerHostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: int64(memory),
		},
		Binds:        []string{path},
		PortBindings: portMap,
	}

	return containerHostConfig
}

func (c *DockerContainer) DataPath() string {
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	return "data/" + c.ID
}
